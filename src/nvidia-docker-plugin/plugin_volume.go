// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/NVIDIA/nvidia-docker/src/nvidia"
)

var (
	ErrVolumeBadFormat   = errors.New("bad volume format")
	ErrVolumeUnsupported = errors.New("unsupported volume")
	ErrVolumeNotFound    = errors.New("no such volume")
	ErrVolumeVersion     = errors.New("invalid volume version")
)

type pluginVolume struct{}

func (p *pluginVolume) implement() string { return "VolumeDriver" }

func (p *pluginVolume) register(api *PluginAPI) {
	prefix := "/" + p.implement()

	api.Handle("POST", prefix+".Create", p.create)
	api.Handle("POST", prefix+".Remove", p.remove)
	api.Handle("POST", prefix+".Mount", p.mount)
	api.Handle("POST", prefix+".Unmount", p.unmount)
	api.Handle("POST", prefix+".Path", p.path)
	api.Handle("POST", prefix+".Get", p.get)
	api.Handle("POST", prefix+".List", p.list)
	api.Handle("POST", prefix+".Capabilities", p.capabilities)
}

func fmtError(err error, vol string) *string {
	s := fmt.Sprintf("%v: %s", err, vol)
	return &s
}

func getVolume(name string) (*nvidia.Volume, string, error) {
	re := regexp.MustCompile("^([a-zA-Z0-9_.-]+)_([0-9.]+)$")
	m := re.FindStringSubmatch(name)
	if len(m) != 3 {
		return nil, "", ErrVolumeBadFormat
	}
	volume, version := Volumes[m[1]], m[2]
	if volume == nil {
		return nil, "", ErrVolumeUnsupported
	}
	return volume, version, nil
}

func (p *pluginVolume) create(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received create request for volume '%s'\n", q.Name)

	volume, version, err := getVolume(q.Name)
	if err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}
	// The volume version requested needs to match the volume version in cache
	if version != volume.Version {
		r.Err = fmtError(ErrVolumeVersion, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}
	ok, err := volume.Exists()
	assert(err)
	if !ok {
		assert(volume.Create(nvidia.LinkStrategy{}))
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) remove(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received remove request for volume '%s'\n", q.Name)

	volume, version, err := getVolume(q.Name)
	if err != nil {
		r.Err = fmtError(err, q.Name)
	} else {
		assert(volume.Remove(version))
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) mount(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Mountpoint, Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received mount request for volume '%s'\n", q.Name)

	volume, version, err := getVolume(q.Name)
	if err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}
	ok, err := volume.Exists(version)
	assert(err)
	if !ok {
		r.Err = fmtError(ErrVolumeNotFound, q.Name)
	} else {
		p := path.Join(volume.Path, version)
		r.Mountpoint = &p
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) unmount(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received unmount request for volume '%s'\n", q.Name)

	_, _, err := getVolume(q.Name)
	if err != nil {
		r.Err = fmtError(err, q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) path(resp http.ResponseWriter, req *http.Request) {
	p.mount(resp, req)
}

func (p *pluginVolume) get(resp http.ResponseWriter, req *http.Request) {
	type Volume struct{ Name, Mountpoint string }

	var q struct{ Name string }
	var r struct {
		Volume *Volume
		Err    *string
	}

	assert(json.NewDecoder(req.Body).Decode(&q))

	volume, version, err := getVolume(q.Name)
	if err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}
	ok, err := volume.Exists(version)
	assert(err)
	if !ok {
		r.Err = fmtError(ErrVolumeNotFound, q.Name)
	} else {
		r.Volume = &Volume{
			Name:       q.Name,
			Mountpoint: path.Join(volume.Path, version),
		}
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) list(resp http.ResponseWriter, req *http.Request) {
	type Volume struct{ Name, Mountpoint string }

	var r struct {
		Volumes []Volume
		Err     *string
	}

	for _, vol := range Volumes {
		versions, err := vol.ListVersions()
		assert(err)
		for _, v := range versions {
			r.Volumes = append(r.Volumes, Volume{
				Name:       fmt.Sprintf("%s_%s", vol.Name, v),
				Mountpoint: path.Join(vol.Path, v),
			})
		}
	}

	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) capabilities(resp http.ResponseWriter, req *http.Request) {
	type Capabilities struct{ Scope string }
	var r struct {
		Capabilities Capabilities
	}

	r.Capabilities = Capabilities{
		Scope: "local",
	}

	assert(json.NewEncoder(resp).Encode(r))
}
