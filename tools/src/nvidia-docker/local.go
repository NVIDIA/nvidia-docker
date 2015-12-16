// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"docker"
	"nvidia"
)

const pluginName = "nvidia"

func CreateLocalVolumes() error {
	vols, err := nvidia.LookupVolumes("")
	if err != nil {
		return err
	}

	for _, v := range vols {
		n := fmt.Sprintf("%s_%s", pluginName, v.Name)
		if err := docker.CreateVolume(n); err != nil {
			return err
		}
		path, err := docker.InspectVolume(n)
		if err != nil {
			return err
		}
		if err := volumeEmpty(n, path); err != nil {
			return err
		}
		if err := v.CreateAt(path); err != nil {
			return err
		}
		fmt.Println(n)
	}
	return nil
}

func volumeEmpty(vol, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF {
		return nil
	}
	return fmt.Errorf("volume %s already exists and is not empty", vol)
}

func GenerateLocalArgs(image string, vols []string) ([]string, error) {
	cv, err := nvidia.GetCUDAVersion()
	if err != nil {
		return nil, err
	}
	if err := cudaSupported(image, cv); err != nil {
		return nil, err
	}

	d, err := devicesArgs()
	if err != nil {
		return nil, err
	}
	v, err := volumesArgs(vols)
	if err != nil {
		return nil, err
	}
	return append(d, v...), nil
}

func devicesArgs() ([]string, error) {
	args := []string{"--device=/dev/nvidiactl", "--device=/dev/nvidia-uvm"}

	// FIXME avoid looking up every devices
	devs, err := nvidia.LookupDevices()
	if err != nil {
		return nil, err
	}

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

	for _, vol := range nvidia.Volumes {
		for _, v := range vols {
			if v == vol.Name {
				// Check if the volume exists locally otherwise fallback to using the plugin
				n := fmt.Sprintf("%s_%s", pluginName, v)
				if _, err := docker.InspectVolume(n); err == nil {
					args = append(args, fmt.Sprintf("--volume=%s:%s", n, vol.Mountpoint))
				} else {
					args = append(args, fmt.Sprintf("--volume-driver=%s", pluginName))
					args = append(args, fmt.Sprintf("--volume=%s:%s", v, vol.Mountpoint))
				}
				break
			}
		}
	}
	return args, nil
}
