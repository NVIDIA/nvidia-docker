// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package nvml

// #cgo LDFLAGS: -ldl -Wl,--unresolved-symbols=ignore-in-object-files
// #include "nvml_dl.h"
import "C"

import (
	"errors"
	"fmt"
)

const (
	szDriver   = C.NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE
	szName     = C.NVML_DEVICE_NAME_BUFFER_SIZE
	szUUID     = C.NVML_DEVICE_UUID_BUFFER_SIZE
	szProcs    = 32
	szProcName = 64
)

type handle struct{ dev C.nvmlDevice_t }

func uintPtr(c C.uint) *uint {
	i := uint(c)
	return &i
}

func uint64Ptr(c C.ulonglong) *uint64 {
	i := uint64(c)
	return &i
}

func stringPtr(c *C.char) *string {
	s := C.GoString(c)
	return &s
}

func errorString(ret C.nvmlReturn_t) error {
	if ret == C.NVML_SUCCESS {
		return nil
	}
	err := C.GoString(C.nvmlErrorString(ret))
	return fmt.Errorf("nvml: %v", err)
}

func init_() error {
	r := C.nvmlInit_dl()
	if r == C.NVML_ERROR_LIBRARY_NOT_FOUND {
		return errors.New("could not load NVML library")
	}
	return errorString(r)
}

func shutdown() error {
	return errorString(C.nvmlShutdown_dl())
}

func systemGetDriverVersion() (string, error) {
	var driver [szDriver]C.char

	r := C.nvmlSystemGetDriverVersion(&driver[0], szDriver)
	return C.GoString(&driver[0]), errorString(r)
}

func systemGetProcessName(pid uint) (string, error) {
	var proc [szProcName]C.char

	r := C.nvmlSystemGetProcessName(C.uint(pid), &proc[0], szProcName)
	return C.GoString(&proc[0]), errorString(r)
}

func deviceGetCount() (uint, error) {
	var n C.uint

	r := C.nvmlDeviceGetCount(&n)
	return uint(n), errorString(r)
}

func deviceGetHandleByIndex(idx uint) (handle, error) {
	var dev C.nvmlDevice_t

	r := C.nvmlDeviceGetHandleByIndex(C.uint(idx), &dev)
	return handle{dev}, errorString(r)
}

func deviceGetTopologyCommonAncestor(h1, h2 handle) (*uint, error) {
	var level C.nvmlGpuTopologyLevel_t

	r := C.nvmlDeviceGetTopologyCommonAncestor_dl(h1.dev, h2.dev, &level)
	if r == C.NVML_ERROR_FUNCTION_NOT_FOUND || r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(C.uint(level)), errorString(r)
}

func (h handle) deviceGetName() (*string, error) {
	var name [szName]C.char

	r := C.nvmlDeviceGetName(h.dev, &name[0], szName)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return stringPtr(&name[0]), errorString(r)
}

func (h handle) deviceGetUUID() (*string, error) {
	var uuid [szUUID]C.char

	r := C.nvmlDeviceGetUUID(h.dev, &uuid[0], szUUID)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return stringPtr(&uuid[0]), errorString(r)
}

func (h handle) deviceGetPciInfo() (*string, error) {
	var pci C.nvmlPciInfo_t

	r := C.nvmlDeviceGetPciInfo(h.dev, &pci)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return stringPtr(&pci.busId[0]), errorString(r)
}

func (h handle) deviceGetMinorNumber() (*uint, error) {
	var minor C.uint

	r := C.nvmlDeviceGetMinorNumber(h.dev, &minor)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(minor), errorString(r)
}

func (h handle) deviceGetBAR1MemoryInfo() (*uint64, *uint64, error) {
	var bar1 C.nvmlBAR1Memory_t

	r := C.nvmlDeviceGetBAR1MemoryInfo(h.dev, &bar1)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil, nil
	}
	return uint64Ptr(bar1.bar1Total), uint64Ptr(bar1.bar1Used), errorString(r)
}

func (h handle) deviceGetPowerManagementLimit() (*uint, error) {
	var power C.uint

	r := C.nvmlDeviceGetPowerManagementLimit(h.dev, &power)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(power), errorString(r)
}

func (h handle) deviceGetMaxClockInfo() (*uint, *uint, error) {
	var sm, mem C.uint

	r := C.nvmlDeviceGetMaxClockInfo(h.dev, C.NVML_CLOCK_SM, &sm)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil, nil
	}
	if r == C.NVML_SUCCESS {
		r = C.nvmlDeviceGetMaxClockInfo(h.dev, C.NVML_CLOCK_MEM, &mem)
	}
	return uintPtr(sm), uintPtr(mem), errorString(r)
}

func (h handle) deviceGetMaxPcieLinkGeneration() (*uint, error) {
	var link C.uint

	r := C.nvmlDeviceGetMaxPcieLinkGeneration(h.dev, &link)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(link), errorString(r)
}

func (h handle) deviceGetMaxPcieLinkWidth() (*uint, error) {
	var width C.uint

	r := C.nvmlDeviceGetMaxPcieLinkWidth(h.dev, &width)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(width), errorString(r)
}

func (h handle) deviceGetPowerUsage() (*uint, error) {
	var power C.uint

	r := C.nvmlDeviceGetPowerUsage(h.dev, &power)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(power), errorString(r)
}

func (h handle) deviceGetTemperature() (*uint, error) {
	var temp C.uint

	r := C.nvmlDeviceGetTemperature(h.dev, C.NVML_TEMPERATURE_GPU, &temp)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(temp), errorString(r)
}

func (h handle) deviceGetUtilizationRates() (*uint, *uint, error) {
	var usage C.nvmlUtilization_t

	r := C.nvmlDeviceGetUtilizationRates(h.dev, &usage)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil, nil
	}
	return uintPtr(usage.gpu), uintPtr(usage.memory), errorString(r)
}

func (h handle) deviceGetEncoderUtilization() (*uint, error) {
	var usage, sampling C.uint

	r := C.nvmlDeviceGetEncoderUtilization(h.dev, &usage, &sampling)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(usage), errorString(r)
}

func (h handle) deviceGetDecoderUtilization() (*uint, error) {
	var usage, sampling C.uint

	r := C.nvmlDeviceGetDecoderUtilization(h.dev, &usage, &sampling)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uintPtr(usage), errorString(r)
}

func (h handle) deviceGetMemoryInfo() (*uint64, error) {
	var mem C.nvmlMemory_t

	r := C.nvmlDeviceGetMemoryInfo(h.dev, &mem)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil
	}
	return uint64Ptr(mem.used), errorString(r)
}

func (h handle) deviceGetClockInfo() (*uint, *uint, error) {
	var sm, mem C.uint

	r := C.nvmlDeviceGetClockInfo(h.dev, C.NVML_CLOCK_SM, &sm)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil, nil
	}
	if r == C.NVML_SUCCESS {
		r = C.nvmlDeviceGetClockInfo(h.dev, C.NVML_CLOCK_MEM, &mem)
	}
	return uintPtr(sm), uintPtr(mem), errorString(r)
}

func (h handle) deviceGetMemoryErrorCounter() (*uint64, *uint64, *uint64, error) {
	var l1, l2, mem C.ulonglong

	r := C.nvmlDeviceGetMemoryErrorCounter(h.dev, C.NVML_MEMORY_ERROR_TYPE_UNCORRECTED,
		C.NVML_VOLATILE_ECC, C.NVML_MEMORY_LOCATION_L1_CACHE, &l1)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil, nil, nil
	}
	if r == C.NVML_SUCCESS {
		r = C.nvmlDeviceGetMemoryErrorCounter(h.dev, C.NVML_MEMORY_ERROR_TYPE_UNCORRECTED,
			C.NVML_VOLATILE_ECC, C.NVML_MEMORY_LOCATION_L2_CACHE, &l2)
	}
	if r == C.NVML_SUCCESS {
		r = C.nvmlDeviceGetMemoryErrorCounter(h.dev, C.NVML_MEMORY_ERROR_TYPE_UNCORRECTED,
			C.NVML_VOLATILE_ECC, C.NVML_MEMORY_LOCATION_DEVICE_MEMORY, &mem)
	}
	return uint64Ptr(l1), uint64Ptr(l2), uint64Ptr(mem), errorString(r)
}

func (h handle) deviceGetPcieThroughput() (*uint, *uint, error) {
	var rx, tx C.uint

	r := C.nvmlDeviceGetPcieThroughput(h.dev, C.NVML_PCIE_UTIL_RX_BYTES, &rx)
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil, nil
	}
	if r == C.NVML_SUCCESS {
		r = C.nvmlDeviceGetPcieThroughput(h.dev, C.NVML_PCIE_UTIL_TX_BYTES, &tx)
	}
	return uintPtr(rx), uintPtr(tx), errorString(r)
}

func (h handle) deviceGetComputeRunningProcesses() ([]uint, []uint64, error) {
	var procs [szProcs]C.nvmlProcessInfo_t
	var count = C.uint(szProcs)

	r := C.nvmlDeviceGetComputeRunningProcesses(h.dev, &count, &procs[0])
	if r == C.NVML_ERROR_NOT_SUPPORTED {
		return nil, nil, nil
	}
	n := int(count)
	pids := make([]uint, n)
	mems := make([]uint64, n)
	for i := 0; i < n; i++ {
		pids[i] = uint(procs[i].pid)
		mems[i] = uint64(procs[i].usedGpuMemory)
	}
	return pids, mems, errorString(r)
}
