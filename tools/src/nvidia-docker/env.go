// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"fmt"
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

func parseAddr(addr string) (host, sport, hport string) {
	re := regexp.MustCompile("^(\\[[0-9a-f.:]+\\]|[0-9A-Za-z.\\-_]+)?(:\\d+)?:(\\d+)?$")

	host, sport, hport = "localhost", "22", "3476"
	if addr == "" {
		return
	}
	m := re.FindStringSubmatch(addr)
	if m == nil {
		return "", "", ""
	}
	if m[1] != "" {
		host = m[1]
	}
	if m[2] != "" {
		sport = m[2][1:]
	}
	if m[3] != "" {
		hport = m[3]
	}
	return
}

func getHost() (u *url.URL) {
	var dhost bool
	var err error

	env := os.Getenv(envNVHost)
	if env == "" {
		env = os.Getenv(envDockerHost)
		dhost = true
	}
	if env == "" {
		return nil
	}

	if ok, _ := regexp.MatchString("^(unix|tcp|http|ssh)://", env); !ok {
		env = "tcp://" + env
	}
	u, err = url.Parse(env)
	if err != nil {
		return nil
	}
	host, sport, hport := parseAddr(u.Host)
	if host == "" {
		return nil
	}

	switch u.Scheme {
	case "tcp":
		u.Scheme = "http"
		fallthrough
	case "http":
		if dhost {
			hport = "3476"
		}
		u.Host = fmt.Sprintf("%s:%s", host, hport)
		return
	case "ssh":
		u.Host = fmt.Sprintf("%s:%s", host, sport)
		u.Opaque = fmt.Sprintf("localhost:%s", hport)
		if u.User == nil {
			u.User = url.UserPassword(os.Getenv("USER"), "")
		}
		return
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
