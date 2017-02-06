// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"fmt"
	"strings"

	"github.com/NVIDIA/nvidia-docker/src/docker"
)

const (
	labelCUDAVersion   = "com.nvidia.cuda.version"
	labelVolumesNeeded = "com.nvidia.volumes.needed"
)

func VolumesNeeded(image string) ([]string, error) {
	ok, err := docker.ImageExists(image)
	if err != nil {
		return nil, err
	}
	if !ok {
		if err = docker.ImagePull(image); err != nil {
			return nil, err
		}
	}

	label, err := docker.Label(image, labelVolumesNeeded)
	if err != nil {
		return nil, err
	}
	if label == "" {
		return nil, nil
	}
	return strings.Split(label, " "), nil
}

func cudaSupported(image, version string) error {
	var vmaj, vmin int
	var lmaj, lmin int

	label, err := docker.Label(image, labelCUDAVersion)
	if err != nil {
		return err
	}
	if label == "" {
		return nil
	}
	if _, err := fmt.Sscanf(version, "%d.%d", &vmaj, &vmin); err != nil {
		return err
	}
	if _, err := fmt.Sscanf(label, "%d.%d", &lmaj, &lmin); err != nil {
		return err
	}
	if lmaj > vmaj || (lmaj == vmaj && lmin > vmin) {
		return fmt.Errorf("unsupported CUDA version: driver %s < image %s", version, label)
	}
	return nil
}
