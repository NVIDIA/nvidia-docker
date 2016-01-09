# Copyright (c) 2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:
.PHONY: all

#################### NVIDIA Samples ####################

all:
	for name in ${CUDA_SAMPLES}; do \
	    $(NV_DOCKER) build -t sample:$$name $$name ; \
	done
