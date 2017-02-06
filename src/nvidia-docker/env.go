// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/NVIDIA/nvidia-docker/src/docker"
)

const (
	envDockerHost = "DOCKER_HOST"
	envNVDocker   = "NV_DOCKER"
	envNVHost     = "NV_HOST"
	envNVGPU      = "NV_GPU"
)

var ErrInvalidURI = errors.New("invalid remote host URI")

func LoadEnvironment() (err error) {
	Host, err = getHost()
	if err != nil {
		return
	}

	GPU = getGPU()
	cmd := getDocker()
	docker.SetCommand(cmd...)
	return
}

func parseAddr(addr string) (host, sport, hport string) {
	re := regexp.MustCompile(`^(\[[0-9a-f.:]+\]|[0-9A-Za-z.\-_]+)?(:\d+)?:(\d+)?$`)

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

func getHost() (*url.URL, error) {
	var env string

	nvhost := os.Getenv(envNVHost)
	dhost := os.Getenv(envDockerHost)

	if nvhost != "" {
		env = nvhost
	} else if dhost != "" {
		env = dhost
	} else {
		return nil, nil
	}

	if nvhost != "" && dhost == "" {
		log.Printf("Warning: %s is set but %s is not\n", envNVHost, envDockerHost)
	}

	if ok, _ := regexp.MatchString("^[a-z0-9+.-]+://", env); !ok {
		env = "tcp://" + env
	}
	uri, err := url.Parse(env)
	if err != nil {
		return nil, ErrInvalidURI
	}
	if uri.Scheme == "unix" {
		return nil, nil
	}

	host, sport, hport := parseAddr(uri.Host)
	if host == "" {
		return nil, ErrInvalidURI
	}

	switch uri.Scheme {
	case "tcp":
		uri.Scheme = "http"
		fallthrough
	case "http":
		if nvhost == "" && dhost != "" {
			hport = "3476"
		}
		uri.Host = fmt.Sprintf("%s:%s", host, hport)
		return uri, nil
	case "ssh":
		uri.Host = fmt.Sprintf("%s:%s", host, sport)
		uri.Opaque = fmt.Sprintf("localhost:%s", hport)
		if uri.User == nil {
			uri.User = url.UserPassword(os.Getenv("USER"), "")
		}
		return uri, nil
	}

	return nil, ErrInvalidURI
}

func getGPU() []string {
	return strings.FieldsFunc(os.Getenv(envNVGPU), func(c rune) bool {
		return c == ' ' || c == ','
	})
}

func getDocker() []string {
	return strings.Fields(os.Getenv(envNVDocker))
}
