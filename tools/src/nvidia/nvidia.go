// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package nvidia

import (
	"errors"
	"os"
	"os/exec"

	"github.com/NVIDIA/nvidia-docker/tools/src/cuda"
	"github.com/NVIDIA/nvidia-docker/tools/src/nvml"
)

const (
	DockerPlugin = "nvidia-docker"
	DeviceCtl    = "/dev/nvidiactl"
	DeviceUVM    = "/dev/nvidia-uvm"
)

func Init() error {
	if err := os.Setenv("CUDA_CACHE_DISABLE", "1"); err != nil {
		return err
	}
	if err := os.Unsetenv("CUDA_VISIBLE_DEVICES"); err != nil {
		return err
	}
	return nvml.Init()
}

func Shutdown() error {
	return nvml.Shutdown()
}

func LoadUVM() error {
	if _, err := os.Stat(DeviceUVM); err == nil {
		return nil
	}
	if exec.Command("nvidia-modprobe", "-u", "-c=0").Run() != nil {
		return errors.New("Could not load UVM kernel module")
	}
	return nil
}

func GetDriverVersion() (string, error) {
	return nvml.GetDriverVersion()
}

func GetCUDAVersion() (string, error) {
	return cuda.GetDriverVersion()
}
