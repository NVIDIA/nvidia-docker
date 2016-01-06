# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:

all: latest

#################### NVIDIA Caffe ####################

latest: $(CURDIR)/Dockerfile
	make -C $(CURDIR)/../cuda 7.0-cudnn4-runtime
	$(NV_DOCKER) build -t caffe $(CURDIR)
