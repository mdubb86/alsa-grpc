package main

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

type ControlInfo struct {
	id       int32
	name     string
	iface    string
	datatype string
	access   string
	numvals  int32
	min      int32
	max      int32
	step     int32
	values   []int32
	percents []int32
}

func parsePair(s string, key string) (string, error) {
	bits := strings.Split(s, "=")
	if len(bits) != 2 {
		return "", fmt.Errorf("invalid pair %s (expected key=value)", s)
	} else if bits[0] != key {
		return "", fmt.Errorf("invalid key in %s (expected %s)", s, key)
	}
	val := strings.TrimPrefix(strings.TrimSuffix(bits[1], "'"), "'")
	return val, nil
}

func parseInt(s string, key string) (int32, error) {
	s, err := parsePair(s, key)
	if err != nil {
		return -1, err
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}

	return int32(i), nil
}

// [numid=2,iface=MIXER,name='Channels'
func parseHeader(ctrl *ControlInfo, line string) error {
	bits := strings.Split(line, ",")

	id, err := parseInt(bits[0], "numid")
	if err != nil {
		return err
	}
	ctrl.id = id

	iface, err := parsePair(bits[1], "iface")
	if err != nil {
		return err
	}
	ctrl.iface = iface

	name, err := parsePair(bits[2], "name")
	if err != nil {
		return err
	}
	ctrl.name = name

	return nil
}

//; type=INTEGER,access=rw---R--,values=2,min=0,max=248,step=0
func parseAttrs(ctrl *ControlInfo, line string) error {
	bits := strings.Split(strings.TrimPrefix(line, "  ; "), ",")

	datatype, err := parsePair(bits[0], "type")
	if err != nil {
		return err
	}
	ctrl.datatype = datatype

	access, err := parsePair(bits[1], "access")
	if err != nil {
		return err
	}
	ctrl.access = access

	numvals, err := parseInt(bits[2], "values")
	if err != nil {
		return err
	}
	ctrl.numvals = numvals

	min, err := parseInt(bits[3], "min")
	if err != nil {
		return err
	}
	ctrl.min = min

	max, err := parseInt(bits[4], "max")
	if err != nil {
		return err
	}
	ctrl.max = max

	step, err := parseInt(bits[5], "step")
	if err != nil {
		return err
	}
	ctrl.step = step

	return nil
}

//: values=181,181
func parseValues(ctrl *ControlInfo, line string) error {
	s, err := parsePair(strings.TrimPrefix(line, "  : "), "values")
	if err != nil {
		return err
	}

	for _, v := range strings.Split(s, ",") {
		num, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		ctrl.values = append(ctrl.values, int32(num))
	}

	return nil
}

func parseControl(lines []string) (*ControlInfo, error) {
	if len(lines) < 3 {
		return nil, fmt.Errorf("control info must be at least 3 lines (got %d lines)", len(lines))
	}

	ctrl := ControlInfo{}
	if err := parseHeader(&ctrl, lines[0]); err != nil {
		return nil, err
	}
	if err := parseAttrs(&ctrl, lines[1]); err != nil {
		return nil, err
	}
	if err := parseValues(&ctrl, lines[2]); err != nil {
		return nil, err
	}

	valRange := ctrl.max - ctrl.min
	if valRange <= 0 {
		return nil, fmt.Errorf("invalid control range for %s (%d - %d)", ctrl.name,  ctrl.min, ctrl.max)
	}

	for _, v := range ctrl.values {
		percent := int32(math.Round(float64(v - ctrl.min) / float64(valRange) * 100.0))
		ctrl.percents = append(ctrl.percents, percent)
	}

	return &ctrl, nil

	// TODO make these test cases
	//test := `numid=2,iface=MIXER,name='Channels'
	//; type=INTEGER,access=rw---R--,values=2,min=0,max=248,step=0
	//: values=181,181
	//| dBscale-min=-100.00dB,step=0.50dB,mute=1`
	//
	//
	//test2 := `numid=6,iface=MIXER,name='Speaker Playback Volume'
	//; type=INTEGER,access=rw---R--,values=2,min=0,max=43,step=0
	//: values=25,25
	//| dBminmax-min=-45.00dB,max=-2.00dB`
	//
	//parseControl(test)
	//parseControl(test2
}

//amixer -D "hw:0" cset numid=2 58%

func cget(device string, numid int) (ControlInfo, error) {
	output, err := exec.Command("amixer", "-D", device, "cget", fmt.Sprintf("numid=%d", numid)).Output()
	if err != nil {
		return ControlInfo{}, err
	}

	ctrl, err := parseControl(strings.Split(string(output), "\n"))
	return *ctrl, err
}

func cset(device string, numid int32, volumes []int32) (ControlInfo, error) {
	volStrs := make([]string, len(volumes))
	for i, v := range volumes {
		volStrs[i] = fmt.Sprintf("%d%", v)
	}

	output, err := exec.Command("amixer", "-D", device, "cset", fmt.Sprintf("numid=%d", numid),
		strings.Join(volStrs, ",")).Output()

	if err != nil {
		return ControlInfo{}, err
	}

	ctrl, err := parseControl(strings.Split(string(output), "\n"))
	return *ctrl, err
}

func amixerContents(device string) ([]ControlInfo, error) {
	output, err := exec.Command("amixer", "-D", device, "contents").Output()
	if err != nil {
		return nil, err
	}

	var ctrls []ControlInfo
	collect := func(lines []string) error {

		if lines == nil {
			return nil
		}

		ctrl, err := parseControl(lines)
		if err != nil {
			return err
		}
		ctrls = append(ctrls, *ctrl)
		return nil
	}

	var chunk []string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "numid") {
			if err := collect(chunk); err != nil {
				return nil, err
			}
			chunk = nil
		}
		chunk = append(chunk, line)
	}
	if err := collect(chunk); err != nil {
		return nil, err
	}

	return ctrls, nil
}
