// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

var dockerCmd = []string{"docker"}

func SetCommand(cmd ...string) {
	if len(cmd) > 0 {
		dockerCmd = cmd
	}
}

func docker(command string, arg ...string) ([]byte, error) {
	var buf bytes.Buffer

	args := append(append(dockerCmd[1:], command), arg...)
	cmd := exec.Command(dockerCmd[0], args...)
	cmd.Stderr = &buf

	b, err := cmd.Output()
	if err != nil {
		b = bytes.TrimSpace(buf.Bytes())
		b = bytes.TrimPrefix(b, []byte("Error: "))
		if len(b) > 0 {
			return nil, fmt.Errorf("%s", b)
		} else {
			return nil, fmt.Errorf("failed to run docker command")
		}
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
	var env []string

	cmd, err := exec.LookPath(dockerCmd[0])
	if err != nil {
		return err
	}
	args := append(dockerCmd, arg...)

	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "DOCKER_") {
			env = append(env, e)
		}
	}
	return syscall.Exec(cmd, args, env)
}
