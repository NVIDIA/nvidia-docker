// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package nvidia

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NVIDIA/nvidia-docker/src/cuda"
	"github.com/NVIDIA/nvidia-docker/src/nvml"
)

type NVMLDevice nvml.Device
type CUDADevice cuda.Device

type Device struct {
	*NVMLDevice
	*CUDADevice
}

type NVMLDeviceStatus nvml.DeviceStatus

type DeviceStatus struct {
	*NVMLDeviceStatus
}

type LookupStrategy uint

const (
	LookupMinimal LookupStrategy = iota
)

func (d *Device) Status() (*DeviceStatus, error) {
	s, err := (*nvml.Device)(d.NVMLDevice).Status()
	if err != nil {
		return nil, err
	}
	return &DeviceStatus{(*NVMLDeviceStatus)(s)}, nil
}

func LookupDevices(s ...LookupStrategy) (devs []Device, err error) {
	var i uint

	n, err := nvml.GetDeviceCount()
	if err != nil {
		return nil, err
	}
	devs = make([]Device, 0, n)
	if n == 0 {
		return
	}

	if len(s) == 1 && s[0] == LookupMinimal {
		for i = 0; i < n; i++ {
			d, err := nvml.NewDeviceLite(i)
			if err != nil {
				return nil, err
			}
			devs = append(devs, Device{(*NVMLDevice)(d), &CUDADevice{}})
		}
		return
	}

	for i = 0; i < n; i++ {
		nd, err := nvml.NewDevice(i)
		if err != nil {
			return nil, err
		}
		cd, err := cuda.NewDevice(nd.PCI.BusID)
		if err != nil {
			return nil, err
		}
		devs = append(devs, Device{(*NVMLDevice)(nd), (*CUDADevice)(cd)})
	}

	for i = 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			ok, err := cuda.CanAccessPeer(
				(*cuda.Device)(devs[i].CUDADevice),
				(*cuda.Device)(devs[j].CUDADevice),
			)
			if err != nil {
				return nil, err
			}
			if ok {
				l, err := nvml.GetP2PLink(
					(*nvml.Device)(devs[i].NVMLDevice),
					(*nvml.Device)(devs[j].NVMLDevice),
				)
				if err != nil {
					return nil, err
				}
				devs[i].Topology = append(devs[i].Topology, nvml.P2PLink{devs[j].PCI.BusID, l})
				devs[j].Topology = append(devs[j].Topology, nvml.P2PLink{devs[i].PCI.BusID, l})
			}
		}
	}
	return
}

func FilterDevices(devs []Device, ids []string) ([]Device, error) {
	type void struct{}
	set := make(map[int]void)

loop:
	for _, id := range ids {
		if strings.HasPrefix(id, "GPU-") {
			for i := range devs {
				if strings.HasPrefix(devs[i].UUID, id) {
					set[i] = void{}
					continue loop
				}
			}
		} else {
			i, err := strconv.Atoi(id)
			if err == nil && i >= 0 && i < len(devs) {
				set[i] = void{}
				continue loop
			}
		}
		return nil, fmt.Errorf("invalid device: %s", id)
	}

	d := make([]Device, 0, len(set))
	for i := range set {
		d = append(d, devs[i])
	}
	return d, nil
}
