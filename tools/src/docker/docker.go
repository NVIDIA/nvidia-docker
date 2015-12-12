// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
)

var dockerBin = []string{"docker"}

func SetBinary(bin ...string) {
	if len(bin) > 0 {
		dockerBin = bin
	}
}

func docker(command string, arg ...string) ([]byte, error) {
	var buf bytes.Buffer

	args := append(append(dockerBin[1:], command), arg...)
	cmd := exec.Command(dockerBin[0], args...)
	cmd.Stderr = &buf

	b, err := cmd.Output()
	if err != nil {
		b = bytes.TrimSpace(buf.Bytes())
		return nil, fmt.Errorf("%s", b)
	}
	return b, nil
}

func ParseArgs(args []string, cmd ...string) (string, int, error) {
	type void struct{}

	re := regexp.MustCompile("(?m)^\\s*(-[^=]+)=[^{true}{false}].*$")
	flags := make(map[string]void)

	b, err := docker("help", cmd...)
	if err != nil {
		return "", -1, err
	}

	// Build the set of Docker flags taking an option using "docker help"
	for _, m := range re.FindAllSubmatch(b, -1) {
		for _, f := range bytes.Split(m[1], []byte(", ")) {
			flags[string(f)] = void{}
		}
	}
	for i := 0; i < len(args); i++ {
		if args[i][:1] == "-" {
			// Skip the flags and their options
			if _, ok := flags[args[i]]; ok {
				i++
			}
			continue
		}
		// Return the first arg that is not a flag
		return args[i], i, nil
	}
	return "", -1, nil
}

func Label(image, label string) (string, error) {
	format := fmt.Sprintf(`--format='{{index .Config.Labels "%s"}}'`, label)

	b, err := docker("inspect", format, image)
	if err != nil {
		return "", err
	}
	return string(bytes.Trim(b, " \n")), nil
}

func CreateVolume(name string) error {
	_, err := docker("volume", "create", "--name", name)
	return err
}

func InspectVolume(name string) (string, error) {
	var vol []struct{ Name, Driver, Mountpoint string }

	b, err := docker("volume", "inspect", name)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(b, &vol); err != nil {
		return "", err
	}
	return vol[0].Mountpoint, nil
}

func Docker(arg ...string) error {
	cmd, err := exec.LookPath(dockerBin[0])
	if err != nil {
		return err
	}
	args := append(dockerBin, arg...)

	return syscall.Exec(cmd, args, nil)
}
