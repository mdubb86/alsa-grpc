package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

type Monitor struct {
	device string
	updates  chan *ControlInfo
	cmd      *exec.Cmd
}

func (m Monitor) Start() error {
	// Use stdbuf for immediate output (disable buffering on stdout)
	m.cmd = exec.Command("/usr/bin/stdbuf", "-oL", "alsactl", "monitor", m.device)
	stdout, err := m.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := m.cmd.Start(); err != nil {
		return err
	}

	r := regexp.MustCompile(`^node (.*), #(\d+) \((\d+),(\d+),(\d+),(\w+),(\d+)\) (\w+)$`)

	go func() {
		fmt.Println("Starting monitoring")
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			// Parsing format derived from https://github.com/alsa-project/alsa-utils/blob/master/alsactl/monitor.c
			bits := r.FindStringSubmatch(scanner.Text())
			deviceName := bits[1]
			numid, err := strconv.Atoi(bits[2])
			//iface := strconv.Atoi(bits[3])
			//device := strconv.Atoi(bits[4])
			//subdevice := strconv.Atoi(bits[5])
			//name := bits[6]
			//index := strconv.Atoi(bits[7])
			if err != nil {
				// TODO log error
			}

			ctlInfo, err := cget(deviceName, numid)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				m.updates <- &ctlInfo
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
