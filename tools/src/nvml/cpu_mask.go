// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package nvml

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type cpuMask [4]uint64

func (m cpuMask) cpuNode() (uint, error) {
	masks, err := getMasks()
	if err != nil {
		return 0, err
	}
	for i, mask := range masks {
		if cmpMasks(m, mask) {
			return uint(i), nil
		}
	}
	return 0, ErrCPUAffinity
}

var cache []cpuMask

func getMasks() (masks []cpuMask, err error) {
	if cache != nil {
		return cache, nil
	}
	err = filepath.Walk("/sys/devices/system/node",
		func(p string, fi os.FileInfo, err error) error {
			var mask cpuMask

			if err != nil {
				return err
			}
			if path.Base(p) != "cpumap" {
				return nil
			}
			f, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			buf := bytes.Split(f[:len(f)-1], []byte(","))
			if len(buf)/2 != len(mask) {
				return ErrCPUAffinity
			}

			for i := range mask {
				h := hex32ToUint64(buf[i*2])
				l := hex32ToUint64(buf[i*2+1])
				mask[len(mask)-1-i] = (h << 32) | l
			}
			masks = append(masks, mask)
			return nil
		})
	cache = masks
	return
}

func cmpMasks(m1, m2 cpuMask) bool {
	for i := range m1 {
		if m1[i] != m2[i] {
			return false
		}
	}
	return true
}

func hex32ToUint64(b []byte) uint64 {
	var n uint32

	h, _ := hex.DecodeString(string(b))
	binary.Read(bytes.NewBuffer(h), binary.BigEndian, &n)
	return uint64(n)
}
