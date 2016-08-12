# Copyright (c) 2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

DIGITS_LATEST := $(word 1, $(DIGITS_VERSIONS))

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:
.PHONY: all latest $(DIGITS_VERSIONS)

all: $(DIGITS_VERSIONS) latest

#################### NVIDIA DIGITS ####################

latest: $(DIGITS_LATEST)
	$(NV_DOCKER) tag digits:$< digits

3.0: $(CURDIR)/3.0/Dockerfile
	make -C $(CURDIR)/../caffe 0.14
	$(NV_DOCKER) build -t digits:$@ $(CURDIR)/$@

3.3: $(CURDIR)/3.3/Dockerfile
	make -C $(CURDIR)/../caffe 0.14
	$(NV_DOCKER) build -t digits:$@ $(CURDIR)/$@

4.0: $(CURDIR)/4.0/Dockerfile
	make -C $(CURDIR)/../caffe 0.15
	$(NV_DOCKER) build -t digits:$@ $(CURDIR)/$@
