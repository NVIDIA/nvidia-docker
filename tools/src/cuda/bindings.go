// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package cuda

// #cgo LDFLAGS: -lcudart_static -ldl -lrt
// #include <stdlib.h>
// #include <cuda_runtime_api.h>
import "C"

import (
	"fmt"
	"unsafe"
)

type handle struct{ dev C.int }

type deviceProp struct {
	major                      int
	minor                      int
	multiProcessorCount        uint
	ECCEnabled                 bool
	totalGlobalMem             uint
	sharedMemPerMultiprocessor uint
	totalConstMem              uint
	l2CacheSize                uint
	memoryClockRate            uint
	memoryBusWidth             uint
}

func errorString(ret C.cudaError_t) error {
	if ret == C.cudaSuccess {
		return nil
	}
	err := C.GoString(C.cudaGetErrorString(ret))
	return fmt.Errorf("cuda: %v", err)
}

func driverGetVersion() (int, error) {
	var driver C.int

	r := C.cudaDriverGetVersion(&driver)
	return int(driver), errorString(r)
}

func deviceGetByPCIBusId(busid string) (handle, error) {
	var dev C.int

	id := C.CString(busid)
	r := C.cudaDeviceGetByPCIBusId(&dev, id)
	C.free(unsafe.Pointer(id))
	return handle{dev}, errorString(r)
}

func deviceCanAccessPeer(h1, h2 handle) (bool, error) {
	var ok C.int

	r := C.cudaDeviceCanAccessPeer(&ok, h1.dev, h2.dev)
	return (ok != 0), errorString(r)
}

func deviceReset() error {
	return errorString(C.cudaDeviceReset())
}

func (h handle) getDeviceProperties() (*deviceProp, error) {
	var props C.struct_cudaDeviceProp

	r := C.cudaGetDeviceProperties(&props, h.dev)
	p := &deviceProp{
		major:                      int(props.major),
		minor:                      int(props.minor),
		multiProcessorCount:        uint(props.multiProcessorCount),
		ECCEnabled:                 bool(props.ECCEnabled != 0),
		totalGlobalMem:             uint(props.totalGlobalMem),
		sharedMemPerMultiprocessor: uint(props.sharedMemPerMultiprocessor),
		totalConstMem:              uint(props.totalConstMem),
		l2CacheSize:                uint(props.l2CacheSize),
		memoryClockRate:            uint(props.memoryClockRate),
		memoryBusWidth:             uint(props.memoryBusWidth),
	}
	return p, errorString(r)
}
