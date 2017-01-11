# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

CUDA_LATEST := $(word 1, $(CUDA_VERSIONS))

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:
.PHONY: all all-cudnn all-cuda latest $(CUDA_VERSIONS) $(CUDNN_VERSIONS)

all: all-cudnn all-cuda

all-cuda: $(addsuffix -runtime, $(CUDA_VERSIONS)) \
          $(addsuffix -devel, $(CUDA_VERSIONS)) \
          latest

all-cudnn: $(addsuffix -runtime, $(CUDNN_VERSIONS)) \
           $(addsuffix -devel, $(CUDNN_VERSIONS))

#################### CUDA ####################

latest: $(CUDA_LATEST)-devel
	$(NV_DOCKER) tag cuda:$< cuda

%-devel: %-runtime $(CURDIR)/%/devel/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/devel

%-runtime: $(CURDIR)/%/runtime/Dockerfile
	$(NV_DOCKER) build -t cuda:$@ $(CURDIR)/$*/runtime

#################### cuDNN ####################

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
