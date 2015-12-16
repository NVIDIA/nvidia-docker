// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"fmt"
	"strings"

	"docker"
)

const (
	labelCUDAVersion   = "com.nvidia.cuda.version"
	labelVolumesNeeded = "com.nvidia.volumes.needed"
)

func VolumesNeeded(image string) ([]string, error) {
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
		return fmt.Errorf("unsupported CUDA version: %s < %s", label, version)
	}
	return nil
}
