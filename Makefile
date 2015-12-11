# Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.
OS ?= ubuntu

# CUDA versions
ifeq ($(OS), ubuntu)
	CUDA_VERSIONS := 7.5 7.0 6.5
else ifeq ($(OS), centos)
        CUDA_VERSIONS := 7.5 7.0
else
$(error unsupported OS: $(OS))
endif
CUDA_LATEST := $(word 1, $(CUDA_VERSIONS))

# cuDNN versions
ifeq ($(OS), ubuntu)
	CUDNN_VERSIONS := 7.5-cudnn4-devel 7.5-cudnn4-runtime \
			  7.5-cudnn3-devel 7.5-cudnn3-runtime \
	                  7.0-cudnn2-devel 7.0-cudnn2-runtime
endif
CUDNN_DEVEL_LATEST := $(word 1, $(CUDNN_VERSIONS))
CUDNN_RUNTIME_LATEST := $(word 2, $(CUDNN_VERSIONS))

# Building Docker images in parallel will duplicate identical layers.
.NOTPARALLEL:

# By default, build only cuda:latest
default: latest

# CUDA images
latest: devel
	docker tag -f cuda:$< cuda

devel: $(CUDA_LATEST)
	docker tag -f cuda:$< cuda:$@

runtime: $(CUDA_LATEST)-runtime
	docker tag -f cuda:$< cuda:$@

%: %-devel $(OS)/cuda/%
	docker tag -f cuda:$< cuda:$@

%-devel: %-runtime $(OS)/cuda/%/devel/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/devel

%-runtime: $(OS)/cuda/%/runtime/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/runtime

all-cuda: $(CUDA_VERSIONS) latest devel runtime

# cuDNN images
cudnn: cudnn-devel
	docker tag -f cuda:$< cuda:$@

cudnn-devel: $(CUDNN_DEVEL_LATEST)
	docker tag -f cuda:$< cuda:$@

cudnn-runtime: $(CUDNN_RUNTIME_LATEST)
	docker tag -f cuda:$< cuda:$@

# Special rules for specific cuDNN versions
%-cudnn2-devel: %-devel $(OS)/cuda/%/devel/cudnn2/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/devel/cudnn2

%-cudnn2-runtime: %-runtime $(OS)/cuda/%/runtime/cudnn2/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/runtime/cudnn2

%-cudnn3-devel: %-devel $(OS)/cuda/%/devel/cudnn3/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/devel/cudnn3

%-cudnn3-runtime: %-runtime $(OS)/cuda/%/runtime/cudnn3/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/runtime/cudnn3

%-cudnn4-devel: %-devel $(OS)/cuda/%/devel/cudnn4/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/devel/cudnn4

%-cudnn4-runtime: %-runtime $(OS)/cuda/%/runtime/cudnn4/Dockerfile
	docker build -t cuda:$@ $(OS)/cuda/$*/runtime/cudnn4

all-cudnn: $(CUDNN_VERSIONS) cudnn cudnn-devel cudnn-runtime

push: all-cuda all-cudnn
	if [ -z "$(USERNAME)" ]; then \
		echo "Error: USERNAME not set"; \
		exit 1; \
	fi; \
        # Retag all images with the username as a prefix.
	docker images | awk '$$1 == "cuda" { print $$2 }' | xargs -I {} docker tag -f cuda:{} $(USERNAME)/cuda:{}
	docker push $(USERNAME)/cuda
	docker images | awk '$$1 == "$(USERNAME)/cuda" { print $$2 }' | xargs -I {} docker rmi $(USERNAME)/cuda:{}

clean:
	docker rmi -f `docker images -q --filter "label=com.nvidia.cuda.version"`
