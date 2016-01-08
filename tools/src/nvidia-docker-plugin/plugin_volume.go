// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"regexp"
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
}

func errVolumeUnknown(vol string) *string {
	s := "No such volume: " + vol
	return &s
}

func parseVolumeName(volume string) (string, string) {
	re := regexp.MustCompile("^([a-zA-Z0-9_.-]+)_([0-9.]+)$")
	if m := re.FindStringSubmatch(volume); len(m) == 3 {
		return m[1], m[2]
	}
	return "", ""
}

func (p *pluginVolume) create(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received create request for volume '%s'\n", q.Name)

	name, version := parseVolumeName(q.Name)
	if v, ok := Volumes[name]; ok {
		if v.Version != version {
			r.Err = new(string)
			*r.Err = "Invalid volume version: " + version
		} else {
			assert(v.Create())
		}
	} else {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) remove(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received remove request for volume '%s'\n", q.Name)

	name, version := parseVolumeName(q.Name)
	if v, ok := Volumes[name]; ok {
		assert(v.Remove(version))
	} else {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) mount(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Mountpoint, Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))

	name, version := parseVolumeName(q.Name)
	if v, ok := Volumes[name]; ok {
		p := path.Join(v.Path, version)
		r.Mountpoint = &p
	} else {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) unmount(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))

	name, _ := parseVolumeName(q.Name)
	if _, ok := Volumes[name]; !ok {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) path(resp http.ResponseWriter, req *http.Request) {
	p.mount(resp, req)
}
