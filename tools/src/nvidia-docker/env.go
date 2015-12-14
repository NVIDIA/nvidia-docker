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
	envDockerHost  = "DOCKER_HOST"
	envNVDockerBin = "NV_DOCKER_BIN"
	envNVHost      = "NV_HOST"
	envNVGPU       = "NV_GPU"
)

func LoadEnvironment() {
	Host = getHost()
	GPU = getGPU()

	b := getDockerBin()
	docker.SetBinary(b...)
}

func getHost() (host *url.URL) {
	re := regexp.MustCompile("([0-9A-Za-z.:-]+):\\d+")

	u, _ := url.Parse(os.Getenv(envNVHost))
	if (u.Scheme == "http" || u.Scheme == "ssh") && re.MatchString(u.Host) {
		return u
	}
	u, _ = url.Parse(os.Getenv(envDockerHost))
	if u.Scheme == "tcp" && re.MatchString(u.Host) {
		u.Scheme = "http"
		u.Host = re.ReplaceAllString(u.Host, "$1:3476")
		return u
	}
	return nil
}

func getGPU() []string {
	return strings.FieldsFunc(os.Getenv(envNVGPU), func(c rune) bool {
		return c == ' ' || c == ','
	})
}

func getDockerBin() []string {
	return strings.Fields(os.Getenv(envNVDockerBin))
}
