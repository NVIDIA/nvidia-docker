# Copyright (c) 2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:
.PHONY: all latest devel runtime

all: latest devel runtime

#################### OpenCL ####################

latest: devel
	$(NV_DOCKER) tag opencl:$< opencl

devel: $(CURDIR)/devel/Dockerfile
	$(NV_DOCKER) build -t opencl:$@ $(CURDIR)/devel

runtime: $(CURDIR)/runtime/Dockerfile
	$(NV_DOCKER) build -t opencl:$@ $(CURDIR)/runtime
