// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package cuda

import (
	"fmt"
)

type MemoryInfo struct {
	ECC       *bool
	Global    *uint
	Shared    *uint
	Constant  *uint
	L2Cache   *uint
	Bandwidth *uint
}

type Device struct {
	handle

	Family *string
	Arch   *string
	Cores  *uint
	Memory MemoryInfo
}

func archFamily(arch string) *string {
	m := map[string]string{
		"1": "Tesla",
		"2": "Fermi",
		"3": "Kepler",
		"5": "Maxwell",
		"6": "Pascal",
	}

	f, ok := m[arch[:1]]
	if !ok {
		return nil
	}
	return &f
}

func archSMCores(arch string) *uint {
	m := map[string]uint{
		"1.0": 8,   // Tesla Generation (SM 1.0) G80 class
		"1.1": 8,   // Tesla Generation (SM 1.1) G8x G9x class
		"1.2": 8,   // Tesla Generation (SM 1.2) GT21x class
		"1.3": 8,   // Tesla Generation (SM 1.3) GT20x class
		"2.0": 32,  // Fermi Generation (SM 2.0) GF100 GF110 class
		"2.1": 48,  // Fermi Generation (SM 2.1) GF10x GF11x class
		"3.0": 192, // Kepler Generation (SM 3.0) GK10x class
		"3.2": 192, // Kepler Generation (SM 3.2) TK1 class
		"3.5": 192, // Kepler Generation (SM 3.5) GK11x GK20x class
		"3.7": 192, // Kepler Generation (SM 3.7) GK21x class
		"5.0": 128, // Maxwell Generation (SM 5.0) GM10x class
		"5.2": 128, // Maxwell Generation (SM 5.2) GM20x class
		"5.3": 128, // Maxwell Generation (SM 5.3) TX1 class
		"6.0": 64,  // Pascal Generation (SM 6.0) GP100 class
		"6.1": 128, // Pascal Generation (SM 6.1) GP10x class
		"6.2": 128, // Pascal Generation (SM 6.2) GP10x class
	}

	c, ok := m[arch]
	if !ok {
		return nil
	}
	return &c
}

func GetDriverVersion() (string, error) {
	d, err := driverGetVersion()
	return fmt.Sprintf("%d.%d", d/1000, d%100/10), err
}

func NewDevice(busid string) (device *Device, err error) {
	h, err := deviceGetByPCIBusId(busid)
	if err != nil {
		return nil, err
	}
	props, err := h.getDeviceProperties()
	if err != nil {
		return nil, err
	}
	arch := fmt.Sprintf("%d.%d", props.major, props.minor)
	family := archFamily(arch)
	cores := archSMCores(arch)
	bw := 2 * (props.memoryClockRate / 1000) * (props.memoryBusWidth / 8)

	// Destroy the active CUDA context
	if err := deviceReset(); err != nil {
		return nil, err
	}

	device = &Device{
		handle: h,
		Family: family,
		Arch:   &arch,
		Cores:  cores,
		Memory: MemoryInfo{
			ECC:       &props.ECCEnabled,
			Global:    &props.totalGlobalMem,
			Shared:    &props.sharedMemPerMultiprocessor,
			Constant:  &props.totalConstMem,
			L2Cache:   &props.l2CacheSize,
			Bandwidth: &bw, // MB/s
		},
	}
	if cores != nil {
		*device.Cores *= props.multiProcessorCount
	}
	*device.Memory.Global /= 1024 * 1024 // MiB
	*device.Memory.Shared /= 1024        // KiB
	*device.Memory.Constant /= 1024      // KiB
	*device.Memory.L2Cache /= 1024       // KiB
	return
}

func CanAccessPeer(dev1, dev2 *Device) (bool, error) {
	return deviceCanAccessPeer(dev1.handle, dev2.handle)
}
