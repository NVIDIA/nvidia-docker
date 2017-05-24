// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/NVIDIA/nvidia-docker/src/docker"
	"github.com/NVIDIA/nvidia-docker/src/nvidia"
)

var (
	Version string
	Host    *url.URL
	GPU     []string
)

func init() {
	log.SetPrefix(os.Args[0] + " | ")
}

func assert(err error) {
	if err != nil {
		log.Panicln("Error:", err)
	}
}

func exit() {
	if err := recover(); err != nil {
		if _, ok := err.(runtime.Error); ok {
			log.Println(err)
		}
		if os.Getenv("NV_DEBUG") != "" {
			log.Printf("%s", debug.Stack())
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	args := os.Args[1:]
	defer exit()

	assert(LoadEnvironment())

	command, off, err := docker.ParseArgs(args)
	assert(err)

	if command == "container" && off+1 < len(args) {
		command = args[off+1]
		off += 1
	}
	if command != "create" && command != "run" {
		if command == "version" {
			fmt.Printf("NVIDIA Docker: %s\n\n", Version)
		}
		assert(docker.Docker(args...))
	}

	opt, i, err := docker.ParseArgs(args[off+1:], command)
	assert(err)
	off += i + 1

	if (command == "create" || command == "run") && opt != "" {
		vols, err := VolumesNeeded(opt)
		assert(err)

		if vols != nil {
			var nargs []string
			var err error

			if Host != nil {
				nargs, err = GenerateRemoteArgs(opt, vols)
			} else {
				assert(nvidia.LoadUVM())
				assert(nvidia.Init())
				nargs, err = GenerateLocalArgs(opt, vols)
				nvidia.Shutdown()
			}
			assert(err)
			args = append(args[:off], append(nargs, args[off:]...)...)
		}
	}

	assert(docker.Docker(args...))
}
