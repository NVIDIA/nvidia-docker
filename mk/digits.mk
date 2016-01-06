# Copyright (c) 2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:

all: latest

#################### NVIDIA DIGITS ####################

latest: $(CURDIR)/Dockerfile
	make -C $(CURDIR)/../caffe
	$(NV_DOCKER) build -t digits $(CURDIR)
