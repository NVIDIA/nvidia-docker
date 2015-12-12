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

func devicesArgs(devs []nvidia.Device) ([]string, error) {
	args := []string{"--device=/dev/nvidiactl", "--device=/dev/nvidia-uvm"}

	if len(GPU) == 0 {
		for i := range devs {
			args = append(args, fmt.Sprintf("--device=%s", devs[i].Path))
		}
	} else {
		for _, id := range GPU {
			i, err := strconv.Atoi(id)
			if err != nil || i < 0 || i >= len(devs) {
				return nil, fmt.Errorf("invalid device: %s", id)
			}
			args = append(args, fmt.Sprintf("--device=%s", devs[i].Path))
		}
	}
	return args, nil
}

func volumesArgs(vols []string) ([]string, error) {
	args := make([]string, 0, len(vols))

	for i := range nvidia.Volumes {
		vol := &nvidia.Volumes[i]

		for _, v := range vols {
			if v == vol.Name {
				// Check if the volume exists locally otherwise fallback to using the plugin
				lv := fmt.Sprintf("%s_%s", PluginName, v)
				if _, err := docker.InspectVolume(lv); err == nil {
					args = append(args, fmt.Sprintf("--volume=%s:%s", lv, vol.Mountpoint))
				} else {
					args = append(args, fmt.Sprintf("--volume-driver=%s", PluginName))
					args = append(args, fmt.Sprintf("--volume=%s:%s", v, vol.Mountpoint))
				}
				break
			}
		}
	}
	return args, nil
}
