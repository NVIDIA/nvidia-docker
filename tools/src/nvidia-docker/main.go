// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"log"
	"os"
	"strings"

	"docker"
	"nvidia"
)

const (
	EnvVolumesPath = "NV_VOLUMES_PATH"
	EnvDockerBin   = "NV_DOCKER_BIN"
	EnvGPU         = "NV_GPU"
)

var (
	VolumesPath string
	DockerBin   []string
	GPU         []string

	Devices []nvidia.Device
	Volumes nvidia.VolumeMap
)

func init() {
	log.SetPrefix(os.Args[0] + " | ")

	if VolumesPath = os.Getenv(EnvVolumesPath); VolumesPath == "" {
		VolumesPath = "/usr/local/nvidia/volumes"
	}
	GPU = strings.FieldsFunc(os.Getenv(EnvGPU), func(c rune) bool {
		return c == ' ' || c == ','
	})
	bin := strings.Fields(os.Getenv(EnvDockerBin))
	docker.SetBinary(bin...)
}

func assert(err error) {
	if err != nil {
		log.Panicln("Error:", err)
	}
}

func exit() {
	code := 0
	if recover() != nil {
		code = 1
	}
	os.Exit(code)
}

func Setup(image string) []string {
	vols, err := volumesNeeded(image)
	assert(err)
	if vols == nil {
		return nil
	}

	assert(nvidia.Init())
	defer func() { assert(nvidia.Shutdown()) }()

	assert(cudaIsSupported(image))
	assert(nvidia.LoadUVM())
	Devices, err = nvidia.GetDevices()
	assert(err)
	Volumes, err = nvidia.GetVolumes(VolumesPath)
	assert(err)
	for _, v := range Volumes {
		assert(v.Create())
	}

	d, err := devicesArgs()
	assert(err)
	v, err := volumesArgs(vols)
	assert(err)
	return append(d, v...)
}

func main() {
	var image string

	args := os.Args[1:]
	defer exit()

	command, i, err := docker.ParseArgs(args)
	assert(err)
	if command != "" {
		image, i, err = docker.ParseArgs(args[i+1:], command)
		assert(err)
	}
	switch command {
	case "create":
		fallthrough
	case "run":
		if image != "" {
			nvargs := Setup(image)
			args = append(args[:i], append(nvargs, args[i:]...)...)
		}
	default:
	}

	assert(docker.Docker(args...))
}
