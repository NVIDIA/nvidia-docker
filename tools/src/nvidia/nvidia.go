// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package nvidia

import (
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
	return exec.Command("nvidia-modprobe", "-u", "-c=0").Run()
}

func GetDriverVersion() (string, error) {
	return nvml.GetDriverVersion()
}

func GetCUDAVersion() (string, error) {
	return cuda.GetDriverVersion()
}
