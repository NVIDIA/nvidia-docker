// Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"flag"
	"log"
	"os"

	"nvidia"
)

var (
	ListenAddr   string
	VolumePrefix string
	SocketPath   string

	Devices []nvidia.Device
	Volumes nvidia.VolumeMap
)

func init() {
	log.SetPrefix(os.Args[0] + " | ")

	flag.StringVar(&ListenAddr, "l", "localhost:3476", "Server listen address")
	flag.StringVar(&VolumePrefix, "p", "", "Volumes prefix path (default is to use mktemp)")
	flag.StringVar(&SocketPath, "s", "/run/docker/plugins/nvidia.sock", "NVIDIA plugin socket path")
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

	if VolumePrefix == "" {
		log.Println("Creating volumes")
	} else {
		log.Println("Creating volumes at", VolumePrefix)
	}
	Volumes, err = nvidia.GetVolumes(VolumePrefix)
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
