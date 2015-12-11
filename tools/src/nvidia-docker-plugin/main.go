// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"nvidia"
)

var (
	ListenAddr  string
	VolumesPath string
	SocketPath  string

	Devices []nvidia.Device
	Volumes nvidia.VolumeMap
)

func init() {
	log.SetPrefix(os.Args[0] + " | ")

	flag.StringVar(&ListenAddr, "l", "localhost:3476", "Server listen address")
	flag.StringVar(&VolumesPath, "v", "", "Path where to store the volumes (default is to use mktemp)")
	flag.StringVar(&SocketPath, "s", "/run/docker/plugins", "NVIDIA plugin socket path")
}

func assert(err error) {
	if err != nil {
		log.Panicln("Error:", err)
	}
}

func exit() {
	code := 0
	if recover() != nil {
		code = 1
	}
	os.Exit(code)
}

func main() {
	var err error

	flag.Parse()
	defer exit()

	log.Println("Loading NVIDIA management library")
	assert(nvidia.Init())
	defer func() { assert(nvidia.Shutdown()) }()

	log.Println("Loading NVIDIA unified memory")
	assert(nvidia.LoadUVM())

	log.Println("Discovering GPU devices")
	Devices, err = nvidia.GetDevices()
	assert(err)

	if VolumesPath == "" {
		VolumesPath, err = ioutil.TempDir("", "nvidia-volumes-")
		assert(err)
		defer func() { assert(os.RemoveAll(VolumesPath)) }()
	}
	log.Println("Creating volumes at", VolumesPath)
	Volumes, err = nvidia.GetVolumes(VolumesPath)
	assert(err)

	plugin := NewPluginAPI(SocketPath)
	remote := NewRemoteAPI(ListenAddr)

	log.Println("Serving plugin API at", SocketPath)
	log.Println("Serving remote API at", ListenAddr)
	p := plugin.Serve()
	r := remote.Serve()

	join, joined := make(chan int, 2), 0
L:
	for {
		select {
		case <-p:
			remote.Stop()
			p = nil
			join <- 1
		case <-r:
			plugin.Stop()
			r = nil
			join <- 1
		case j := <-join:
			if joined += j; joined == cap(join) {
				break L
			}
		}
	}
	assert(plugin.Error())
	assert(remote.Error())
	log.Println("Successfully terminated")
}
