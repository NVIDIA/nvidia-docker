// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"encoding/json"
	"net/http"
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

func sptr(s string) *string { return &s }

func errVolumeUnknown(vol string) *string {
	return sptr("No such volume: " + vol)
}

func (p *pluginVolume) create(resp http.ResponseWriter, req *http.Request) {
	q := struct{ Name string }{}
	r := struct{ Err *string }{}

	assert(json.NewDecoder(req.Body).Decode(&q))
	if v, ok := Volumes[q.Name]; ok {
		if err := v.Create(); err != nil {
			r.Err = sptr(err.Error())
		}
	} else {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) remove(resp http.ResponseWriter, req *http.Request) {
	q := struct{ Name string }{}
	r := struct{ Err *string }{}

	assert(json.NewDecoder(req.Body).Decode(&q))
	if v, ok := Volumes[q.Name]; ok {
		if err := v.Remove(); err != nil {
			r.Err = sptr(err.Error())
		}
	} else {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) mount(resp http.ResponseWriter, req *http.Request) {
	q := struct{ Name string }{}
	r := struct{ Mountpoint, Err *string }{}

	assert(json.NewDecoder(req.Body).Decode(&q))
	if v, ok := Volumes[q.Name]; ok {
		r.Mountpoint = &v.Path
	} else {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) unmount(resp http.ResponseWriter, req *http.Request) {
	q := struct{ Name string }{}
	r := struct{ Err *string }{}

	assert(json.NewDecoder(req.Body).Decode(&q))
	if _, ok := Volumes[q.Name]; !ok {
		r.Err = errVolumeUnknown(q.Name)
	}
	assert(json.NewEncoder(resp).Encode(r))
}

func (p *pluginVolume) path(resp http.ResponseWriter, req *http.Request) {
	p.mount(resp, req)
}
