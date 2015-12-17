// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"docker"
)

const (
	envDockerHost = "DOCKER_HOST"
	envNVDocker   = "NV_DOCKER"
	envNVHost     = "NV_HOST"
	envNVGPU      = "NV_GPU"
)

func LoadEnvironment() {
	Host = getHost()
	GPU = getGPU()

	cmd := getDocker()
	docker.SetCommand(cmd...)
}

func getHost() (host *url.URL) {
	var err error

	re := regexp.MustCompile("^([0-9A-Za-z.:\\-\\[\\]]+)(:\\d+)$")

	if h := os.Getenv(envNVHost); h != "" {
		host, err = url.Parse(h)
		if err != nil {
			return nil
		}
	} else {
		host, err = url.Parse(os.Getenv(envDockerHost))
		if err != nil {
			return nil
		}
		if host.Scheme == "tcp" {
			host.Scheme = "ssh"
			host.Host = re.ReplaceAllString(host.Host, "$1:3476")
		}
	}
	if re.MatchString(host.Host) {
		switch host.Scheme {
		case "ssh":
			m := re.FindStringSubmatch(host.Host)
			host.Host = m[1]
			if !re.MatchString(host.Host) {
				host.Host += ":22"
			}
			host.Opaque = "localhost" + m[2]
			return
		case "http":
			return
		}
	}
	return nil
}

func getGPU() []string {
	return strings.FieldsFunc(os.Getenv(envNVGPU), func(c rune) bool {
		return c == ' ' || c == ','
	})
}

func getDocker() []string {
	return strings.Fields(os.Getenv(envNVDocker))
}
