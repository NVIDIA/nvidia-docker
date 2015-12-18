# Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:

all: latest

#################### NVIDIA Caffe ####################

latest: $(CURDIR)/Dockerfile
	$(NV_DOCKER) build -t caffe $(CURDIR)
