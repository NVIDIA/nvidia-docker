// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var dockerCmd = []string{"docker"}

func SetCommand(cmd ...string) {
	if len(cmd) > 0 {
		dockerCmd = cmd
	}
}

func docker(stdout bool, command string, arg ...string) (b []byte, err error) {
	var buf bytes.Buffer

	args := append(append(dockerCmd[1:], command), arg...)
	cmd := exec.Command(dockerCmd[0], args...)
	cmd.Stderr = &buf

	if stdout {
		cmd.Stdout = os.Stderr
		err = cmd.Run()
	} else {
		b, err = cmd.Output()
	}
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

// List of boolean options: https://github.com/docker/docker/blob/17.03.x/contrib/completion/bash/docker
var lastSupportedVersion = "17.03"
var booleanFlags = map[string]map[string][]string{
	"1.9": {
		"": []string{"-debug", "D", "-tls", "-tlsverify"}, // global options
		"daemon": []string{"-debug", "D", "-tls", "-tlsverify", // global options
			"-disable-legacy-registry", "-help", "-icc", "-ip-forward",
			"-ip-masq", "-iptables", "-ipv6", "-selinux-enabled", "-userland-proxy"},
		"create": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t"},
		"run": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t", // same as "create"
			"-detach", "d", "-rm", "-sig-proxy"},
	},
	"1.10": {
		"": []string{"-debug", "D", "-tls", "-tlsverify"}, // global options
		"daemon": []string{"-debug", "D", "-tls", "-tlsverify", // global options
			"-disable-legacy-registry", "-help", "-icc", "-ip-forward",
			"-ip-masq", "-iptables", "-ipv6", "-selinux-enabled", "-userland-proxy"},
		"create": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t"},
		"run": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t", // same as "create"
			"-detach", "d", "-rm", "-sig-proxy"},
	},
	"1.11": {
		"": []string{"-debug", "D", "-tls", "-tlsverify"}, // global options
		"daemon": []string{"-debug", "D", "-tls", "-tlsverify", // global options
			"-disable-legacy-registry", "-help", "-icc", "-ip-forward",
			"-ip-masq", "-iptables", "-ipv6", "-raw-logs", "-selinux-enabled", "-userland-proxy"},
		"create": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t"},
		"run": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t", // same as "create"
			"-detach", "d", "-rm", "-sig-proxy"},
	},
	"1.12": {
		"": []string{"-debug", "D", "-tls", "-tlsverify"}, // global options
		"daemon": []string{"-debug", "D", "-tls", "-tlsverify", // global options
			"-disable-legacy-registry", "-help", "-icc", "-ip-forward",
			"-ip-masq", "-iptables", "-ipv6", "-live-restore", "-raw-logs",
			"-selinux-enabled", "-userland-proxy"},
		"create": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t"},
		"run": []string{"-disable-content-trust", "-help", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t", // same as "create"
			"-detach", "d", "-no-healthcheck", "-rm", "-sig-proxy"},
	},
	"1.13": {
		"": []string{"-debug", "D", "-tls", "-tlsverify"}, // global options
		"daemon": []string{"-debug", "D", "-tls", "-tlsverify", // global options
			"-disable-legacy-registry", "-experimental", "-help", "-icc", "-init", "-ip-forward",
			"-ip-masq", "-iptables", "-ipv6", "-live-restore", "-raw-logs",
			"-selinux-enabled", "-userland-proxy"},
		"create": []string{"-disable-content-trust", "-help", "-init", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t"},
		"run": []string{"-disable-content-trust", "-help", "-init", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t", // same as "create"
			"-detach", "d", "-no-healthcheck", "-rm", "-sig-proxy"},
	},
	lastSupportedVersion: {
		"": []string{"-debug", "D", "-tls", "-tlsverify"}, // global options
		"daemon": []string{"-debug", "D", "-tls", "-tlsverify", // global options
			"-disable-legacy-registry", "-experimental", "-help", "-icc", "-init", "-ip-forward",
			"-ip-masq", "-iptables", "-ipv6", "-live-restore", "-raw-logs",
			"-selinux-enabled", "-userland-proxy"},
		"create": []string{"-disable-content-trust", "-help", "-init", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t"},
		"run": []string{"-disable-content-trust", "-help", "-init", "-interactive", "i", "-oom-kill-disable",
			"-privileged", "-publish-all", "P", "-read-only", "-tty", "t", // same as "create"
			"-detach", "d", "-no-healthcheck", "-rm", "-sig-proxy"},
	},
}

func ParseArgs(args []string, cmd ...string) (string, int, error) {
	if len(cmd) == 0 {
		cmd = append(cmd, "")
	}
	version, err := ClientVersion()
	if err != nil {
		return "", -1, err
	}
	vmaj := version[:strings.LastIndex(version, ".")]

	cmdBooleanFlags, ok := booleanFlags[vmaj][cmd[0]]
	if !ok {
		// Docker is newer than supported version: use flags from last version we know.
		cmdBooleanFlags, _ = booleanFlags[lastSupportedVersion][cmd[0]]
	}

	// Build the set of boolean Docker options for this command
	type void struct{}
	flags := make(map[string]void)
	for _, f := range cmdBooleanFlags {
		flags[f] = void{}
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg[0] != '-' || arg == "-" {
			return args[i], i, nil
		}
		// Skip if current flag is in the form --option=value
		// Note: doesn't handle weird commands like `nvidia-docker run -vit=XYZ /tmp:/bar ubuntu`
		if strings.Contains(arg, "=") {
			continue
		}

		arg = arg[1:]
		if arg[0] == '-' {
			// Long option: skip next argument if option is not boolean
			if _, ok := flags[arg]; !ok {
				i++
			}
		} else {
			// Short options: skip next argument if any option is not boolean
			for _, f := range arg {
				if _, ok := flags[string(f)]; !ok {
					i++
					break
				}
			}
		}
	}
	return "", -1, nil
}

func Label(image, label string) (string, error) {
	format := fmt.Sprintf(`--format={{index .Config.Labels "%s"}}`, label)

	b, err := docker(false, "inspect", format, image)
	if err != nil {
		return "", err
	}
	return string(bytes.Trim(b, " \n")), nil
}

func VolumeInspect(name string) (string, error) {
	var vol []struct{ Name, Driver, Mountpoint string }

	b, err := docker(false, "volume", "inspect", name)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(b, &vol); err != nil {
		return "", err
	}
	return vol[0].Mountpoint, nil
}

func ImageExists(image string) (bool, error) {
	_, err := docker(false, "inspect", "--type=image", image)
	if err != nil {
		// We can't know whether the image was missing or if the daemon was unreachable.
		return false, nil
	}

	return true, nil
}

func ImagePull(image string) error {
	_, err := docker(true, "pull", image)
	return err
}

func ClientVersion() (string, error) {
	b, err := docker(false, "version", "--format", "{{.Client.Version}}")
	if err != nil {
		return "", err
	}
	version := string(b)
	var v1, v2, v3 int
	if _, err := fmt.Sscanf(version, "%d.%d.%d", &v1, &v2, &v3); err != nil {
		return "", err
	}
	return version, nil
}

func Docker(arg ...string) error {
	cmd, err := exec.LookPath(dockerCmd[0])
	if err != nil {
		return err
	}
	args := append(dockerCmd, arg...)

	return syscall.Exec(cmd, args, os.Environ())
}
