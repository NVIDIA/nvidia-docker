# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

CAFFE_LATEST := $(word 1, $(CAFFE_VERSIONS))

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:
.PHONY: all latest $(CAFFE_VERSIONS)

all: $(CAFFE_VERSIONS) latest

#################### NVIDIA Caffe ####################

latest: $(CAFFE_LATEST)
	$(NV_DOCKER) tag caffe:$< caffe

0.14: $(CURDIR)/0.14/Dockerfile
	make -C $(CURDIR)/../cuda 7.5-cudnn5-runtime
	$(NV_DOCKER) build -t caffe:$@ $(CURDIR)/$@
