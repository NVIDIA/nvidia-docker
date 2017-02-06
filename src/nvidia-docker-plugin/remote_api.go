// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"net/http"

	"github.com/NVIDIA/nvidia-docker/src/graceful"
)

type restapi interface {
	version() string

	gpuInfo(http.ResponseWriter, *http.Request)
	gpuInfoJSON(http.ResponseWriter, *http.Request)
	gpuStatus(http.ResponseWriter, *http.Request)
	gpuStatusJSON(http.ResponseWriter, *http.Request)
	dockerCLI(http.ResponseWriter, *http.Request)
	dockerCLIJSON(http.ResponseWriter, *http.Request)
	mesosCLI(http.ResponseWriter, *http.Request)
}

type RemoteAPI struct {
	*graceful.HTTPServer

	apis []restapi
}

func NewRemoteAPI(addr string) *RemoteAPI {
	a := &RemoteAPI{
		HTTPServer: graceful.NewHTTPServer("tcp", addr),
	}
	a.register(
		new(remoteV10),
	)
	return a
}

func (a *RemoteAPI) register(apis ...restapi) {
	for i, api := range apis {
		prefix := "/" + api.version()

	handlers:
		a.Handle("GET", prefix+"/gpu/info", api.gpuInfo)
		a.Handle("GET", prefix+"/gpu/info/json", api.gpuInfoJSON)
		a.Handle("GET", prefix+"/gpu/status", api.gpuStatus)
		a.Handle("GET", prefix+"/gpu/status/json", api.gpuStatusJSON)
		a.Handle("GET", prefix+"/docker/cli", api.dockerCLI)
		a.Handle("GET", prefix+"/docker/cli/json", api.dockerCLIJSON)
		a.Handle("GET", prefix+"/mesos/cli", api.mesosCLI)

		if i == len(apis)-1 && prefix != "" {
			prefix = ""
			goto handlers
		}
		a.apis = append(a.apis, api)
	}
}
