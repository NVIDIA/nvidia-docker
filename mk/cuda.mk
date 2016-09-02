# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

CUDA_LATEST := $(word 1, $(CUDA_VERSIONS))
CUDNN_LATEST := $(word 1, $(CUDNN_VERSIONS))

BUILD_ARCH = $(shell uname -m)
ifneq ($(BUILD_ARCH),ppc64le)
    BUILD_ARCH = ''
endif

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:
.PHONY: all all-cudnn all-cuda latest $(CUDA_VERSIONS) $(CUDNN_VERSIONS)

all: all-cudnn all-cuda

all-cuda: $(CUDA_VERSIONS) latest devel runtime

all-cudnn: $(addsuffix -runtime, $(CUDNN_VERSIONS)) \
           $(addsuffix -devel, $(CUDNN_VERSIONS)) \
           cudnn cudnn-devel cudnn-runtime

#################### CUDA ####################

latest: devel
	$(NV_DOCKER) tag cuda:$< cuda

devel: $(CUDA_LATEST)
	$(NV_DOCKER) tag cuda:$< cuda:$@

runtime: $(CUDA_LATEST)-runtime
	$(NV_DOCKER) tag cuda:$< cuda:$@

8.0: 8.0-devel $(CURDIR)/8.0
	$(NV_DOCKER) tag cuda:$< cuda:$@

7.5: 7.5-devel $(CURDIR)/7.5
	$(NV_DOCKER) tag cuda:$< cuda:$@

7.0: 7.0-devel $(CURDIR)/7.0
	$(NV_DOCKER) tag cuda:$< cuda:$@

6.5: 6.5-devel $(CURDIR)/6.5
	$(NV_DOCKER) tag cuda:$< cuda:$@

%-devel: %-runtime $(CURDIR)/%/devel/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/devel

%-runtime: $(CURDIR)/%/runtime/Dockerfile
	$(NV_DOCKER) build -f $(CURDIR)/$*/runtime/Dockerfile.$(BUILD_ARCH) -t cuda:$@ $(CURDIR)/$*/runtime

#################### cuDNN ####################

cudnn: cudnn-devel
	$(NV_DOCKER) tag cuda:$< cuda:$@

cudnn-devel: $(CUDNN_LATEST)-devel
	$(NV_DOCKER) tag cuda:$< cuda:$@

cudnn-runtime: $(CUDNN_LATEST)-runtime
	$(NV_DOCKER) tag cuda:$< cuda:$@

%-cudnn2-devel: %-devel $(CURDIR)/%/devel/cudnn2/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/devel/cudnn2

%-cudnn2-runtime: %-runtime $(CURDIR)/%/runtime/cudnn2/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/runtime/cudnn2

%-cudnn3-devel: %-devel $(CURDIR)/%/devel/cudnn3/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/devel/cudnn3

%-cudnn3-runtime: %-runtime $(CURDIR)/%/runtime/cudnn3/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/runtime/cudnn3

%-cudnn4-devel: %-devel $(CURDIR)/%/devel/cudnn4/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/devel/cudnn4

%-cudnn4-runtime: %-runtime $(CURDIR)/%/runtime/cudnn4/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/runtime/cudnn4

%-cudnn5-devel: %-devel $(CURDIR)/%/devel/cudnn5/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/devel/cudnn5

%-cudnn5-runtime: %-runtime $(CURDIR)/%/runtime/cudnn5/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/runtime/cudnn5
