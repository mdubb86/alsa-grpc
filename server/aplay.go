package main

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// TODO make test cases
//`
//**** List of PLAYBACK Hardware Devices ****
//card 0: HDMI [HDA Intel HDMI], device 3: HDMI 0 [HDMI 0]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//card 0: HDMI [HDA Intel HDMI], device 7: HDMI 1 [HDMI 1]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//card 0: HDMI [HDA Intel HDMI], device 8: HDMI 2 [HDMI 2]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//card 0: HDMI [HDA Intel HDMI], device 9: HDMI 3 [HDMI 3]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//card 0: HDMI [HDA Intel HDMI], device 10: HDMI 4 [HDMI 4]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//card 1: system [iMic USB audio system], device 0: USB Audio [USB Audio]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//card 2: PCH [HDA Intel PCH], device 0: ALC3220 Analog [ALC3220 Analog]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//`
//`
//**** List of PLAYBACK Hardware Devices ****
//card 0: sndrpihifiberry [snd_rpi_hifiberry_amp], device 0: HifiBerry AMP HiFi tas5713.1-001b-0 [HifiBerry AMP HiFi tas5713.1-001b-0]
//  Subdevices: 1/1
//  Subdevice #0: subdevice #0
//`

func aplayList() (map[string]int32, error) {
	output, err := exec.Command("aplay", "-l").Output()
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`^card (\d+): [A-Za-z0-9_ ]+ \[([A-Za-z0-9_ ]+)]`)
	cards := map[string]int32{}
	for _, line := range strings.Split(string(output), "\n") {
		bits := r.FindStringSubmatch(line)
		if len(bits) == 3 {
			cardNum, err := strconv.Atoi(bits[1])
			if err != nil {
				return nil, err
			}
			name := bits[2]
			cards[name] = int32(cardNum)
		}
	}
	return cards, nil
}
