package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

type Monitor struct {
	cardNum  int32
	updates  chan *ControlInfo
	cmd      *exec.Cmd
}

func (m Monitor) Start() error {
	// Use stdbuf for immediate output (disable buffering on stdout)
	m.cmd = exec.Command("/usr/bin/stdbuf", "-oL", "alsactl", "monitor")
	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := m.cmd.Start(); err != nil {
		return err
	}

	// Parsing format derived from https://github.com/alsa-project/alsa-utils/blob/master/alsactl/monitor.c
	// node hw:0, #1 (2,0,0,Master,0) VALUE
	// card 1, #8 (2,0,0,Mic Capture Volume,0) VALUE
	r := regexp.MustCompile(`^(node|card) (hw:)?(\d+), #(\d+) \((\d+),(\d+),(\d+),([A-Za-z0-9\-_ ]+),(\d+)\) (\w+)$`)

	go func() {
		fmt.Println("Starting monitoring")
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			bits := r.FindStringSubmatch(scanner.Text())

			if len(bits) != 11 {
				fmt.Println("Ignoring unrecognized update: ", scanner.Text())
			} else {
				cardNum, _ := strconv.Atoi(bits[3])
				ctrlId, _ := strconv.Atoi(bits[4])
				ctlInfo, err := cget(cardNum, ctrlId)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					m.updates <- &ctlInfo
				}
			}
		}
	}()
	return nil
}

// Code for killing process
//err := cmd.Process.Kill()
//if err != nil {
//    panic(err) // panic as can't kill a process.
//}
