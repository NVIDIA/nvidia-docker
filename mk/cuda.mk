# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

CUDA_LATEST := $(word 1, $(CUDA_VERSIONS))
CUDNN_LATEST := $(word 1, $(CUDNN_VERSIONS))

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:

all: all-cudnn all-cuda

all-cuda: $(CUDA_VERSIONS) latest devel runtime

all-cudnn: $(addsuffix -runtime, $(CUDNN_VERSIONS)) \
           $(addsuffix -devel, $(CUDNN_VERSIONS)) \
           cudnn cudnn-devel cudnn-runtime

#################### CUDA ####################

latest: devel
	$(NV_DOCKER) tag -f cuda:$< cuda

devel: $(CUDA_LATEST)
	$(NV_DOCKER) tag -f cuda:$< cuda:$@

runtime: $(CUDA_LATEST)-runtime
	$(NV_DOCKER) tag -f cuda:$< cuda:$@

%: %-devel $(CURDIR)/%
	$(NV_DOCKER) tag -f cuda:$< cuda:$@

%-devel: %-runtime $(CURDIR)/%/devel/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/devel

%-runtime: $(CURDIR)/%/runtime/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/runtime

#################### cuDNN ####################

cudnn: cudnn-devel
	$(NV_DOCKER) tag -f cuda:$< cuda:$@

cudnn-devel: $(CUDNN_LATEST)-devel
	$(NV_DOCKER) tag -f cuda:$< cuda:$@

cudnn-runtime: $(CUDNN_LATEST)-runtime
	$(NV_DOCKER) tag -f cuda:$< cuda:$@

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
