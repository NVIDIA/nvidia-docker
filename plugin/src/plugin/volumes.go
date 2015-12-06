// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"bufio"
	"bytes"
	"debug/elf"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"ldcache"
)

const (
	tempDir  = "nvidia-volumes-"
	lib32Dir = "lib"
	lib64Dir = "lib64"
)

type volumeDir struct {
	name  string
	files []string
}

type volumeInfo struct {
	Path string
	dirs []volumeDir
}

type Volume struct {
	Name       string
	Mountpoint string
	components []string

	*volumeInfo
}

type VolumeMap map[string]*Volume

var volumes = []Volume{
	{
		"bin",
		"/usr/local/nvidia/bin",
		[]string{
			"nvidia-cuda-mps-control",
			"nvidia-cuda-mps-server",
			"nvidia-debugdump",
			"nvidia-persistenced",
			"nvidia-smi",
		}, nil,
	},
	{
		"cuda",
		"/usr/local/nvidia",
		[]string{
			"libcuda.so",
			"libnvcuvid.so",
			"libnvidia-compiler.so",
			"libnvidia-encode.so",
			"libnvidia-ml.so",
		}, nil,
	},
}

func (v *Volume) Create() (err error) {
	if err = os.MkdirAll(v.Path, 0755); err != nil {
		return
	}
	defer func() {
		if err != nil {
			v.Remove()
		}
	}()

	for _, d := range v.dirs {
		dir := path.Join(v.Path, d.name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		for _, f := range d.files {
			obj, err := elf.Open(f)
			if err != nil {
				return err
			}
			soname, err := obj.DynString(elf.DT_SONAME)
			obj.Close()
			if err != nil {
				return err
			}

			l := path.Join(dir, path.Base(f))
			if err := os.Link(f, l); err != nil {
				return err
			}
			if len(soname) > 0 {
				f = path.Join(v.Mountpoint, d.name, path.Base(f))
				l = path.Join(dir, soname[0])
				if err := os.Symlink(f, l); err != nil &&
					!os.IsExist(err.(*os.LinkError).Err) {
					return err
				}
			}
		}
	}
	return nil
}

func (v *Volume) Remove() error {
	return os.RemoveAll(v.Path)
}

func which(bins ...string) ([]string, error) {
	paths := make([]string, 0, len(bins))

	out, _ := exec.Command("which", bins...).Output()
	r := bufio.NewReader(bytes.NewBuffer(out))
	for {
		p, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if p = strings.TrimSpace(p); !path.IsAbs(p) {
			continue
		}
		path, err := filepath.EvalSymlinks(p)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func GetVolumes(prefix string) (vols VolumeMap, err error) {
	if prefix == "" {
		prefix, err = ioutil.TempDir("", tempDir)
		if err != nil {
			return
		}
	}

	cache, err := ldcache.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := cache.Close(); err == nil {
			err = e
		}
	}()

	vols = make(VolumeMap, len(volumes))

	for i := range volumes {
		vol := &volumes[i]
		vol.volumeInfo = &volumeInfo{
			Path: path.Join(prefix, vol.Name),
		}

		if vol.Name == "bin" {
			bins, err := which(vol.components...)
			if err != nil {
				return nil, err
			}
			vol.dirs = append(vol.dirs, volumeDir{".", bins})
		} else {
			libs32, libs64 := cache.Lookup(vol.components...)
			vol.dirs = append(vol.dirs,
				volumeDir{lib32Dir, libs32},
				volumeDir{lib64Dir, libs64},
			)
		}
		vols[vol.Name] = vol
	}
	return
}
