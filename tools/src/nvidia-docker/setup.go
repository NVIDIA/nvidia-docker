// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"fmt"
	"strconv"
	"strings"

	"docker"
	"nvidia"
)

const (
	labelCUDAVersion   = "com.nvidia.cuda.version"
	labelVolumesNeeded = "com.nvidia.volumes.needed"
)

func cudaIsSupported(image string) error {
	var vmaj, vmin int
	var lmaj, lmin int

	label, err := docker.Label(image, labelCUDAVersion)
	if err != nil {
		return err
	}
	if label == "" {
		return nil
	}
	version, err := nvidia.GetCUDAVersion()
	if err != nil {
		return err
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

func volumesNeeded(image string) ([]string, error) {
	label, err := docker.Label(image, labelVolumesNeeded)
	if err != nil {
		return nil, err
	}
	if label == "" {
		return nil, nil
	}
	return strings.Split(label, " "), nil
}

func devicesArgs() ([]string, error) {
	args := []string{"--device=/dev/nvidiactl", "--device=/dev/nvidia-uvm"}

	if len(GPU) == 0 {
		for i := range Devices {
			args = append(args, fmt.Sprintf("--device=%s", Devices[i].Path))
		}
	} else {
		for _, id := range GPU {
			i, err := strconv.Atoi(id)
			if err != nil || i < 0 || i >= len(Devices) {
				return nil, fmt.Errorf("invalid device: %s", id)
			}
			args = append(args, fmt.Sprintf("--device=%s", Devices[i].Path))
		}
	}
	return args, nil
}

func volumesArgs(needed []string) ([]string, error) {
	args := make([]string, 0, len(needed))

	for _, n := range needed {
		v, ok := Volumes[n]
		if !ok {
			return nil, fmt.Errorf("invalid volume: %s", n)
		}
		args = append(args, fmt.Sprintf("--volume=%s:%s", v.Path, v.Mountpoint))
	}
	return args, nil
}
