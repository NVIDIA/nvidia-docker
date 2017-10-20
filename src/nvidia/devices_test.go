// Copyright (c) 2017, NVIDIA CORPORATION. All rights reserved.

package nvidia

import (
	"testing"
)

func TestFilterDevices(t *testing.T) {
	nd1 := NVMLDevice{UUID: "GPU-1"}
	d1 := Device{(*NVMLDevice)(&nd1), (*CUDADevice)(nil)}

	nd2 := NVMLDevice{UUID: "GPU-2"}
	d2 := Device{(*NVMLDevice)(&nd2), (*CUDADevice)(nil)}

	devs, err := FilterDevices([]Device{d1, d2}, []string{"GPU-1"})
	if err != nil {
		t.Error(err)
	}

	if len(devs) != 1 {
		t.Errorf("unexpected number of devices: %d", len(devs))
	}
	d := devs[0]
	if d.UUID != "GPU-1" {
		t.Errorf("unexpected uuid: %s", d.UUID)
	}
}
