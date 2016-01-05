# Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.
OS ?= ubuntu
DOCKER_BIN ?= docker

# CUDA versions
ifeq ($(OS), ubuntu)
	CUDA_VERSIONS := 7.5 7.0 6.5
else ifeq ($(OS), centos-7)
        CUDA_VERSIONS := 7.5 7.0
else ifeq ($(OS), centos-6)
        CUDA_VERSIONS := 7.5
else
$(error unsupported OS: $(OS))
endif
CUDA_LATEST := $(word 1, $(CUDA_VERSIONS))

# cuDNN versions
ifeq ($(OS), ubuntu)
	CUDNN_VERSIONS := 7.5-cudnn4-devel 7.5-cudnn4-runtime \
			  7.5-cudnn3-devel 7.5-cudnn3-runtime \
			  7.0-cudnn4-devel 7.0-cudnn4-runtime \
			  7.0-cudnn3-devel 7.0-cudnn3-runtime \
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
	$(DOCKER_BIN) tag -f cuda:$< cuda

devel: $(CUDA_LATEST)
	$(DOCKER_BIN) tag -f cuda:$< cuda:$@

runtime: $(CUDA_LATEST)-runtime
	$(DOCKER_BIN) tag -f cuda:$< cuda:$@

%: %-devel $(OS)/cuda/%
	$(DOCKER_BIN) tag -f cuda:$< cuda:$@

%-devel: %-runtime $(OS)/cuda/%/devel/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/devel

%-runtime: $(OS)/cuda/%/runtime/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/runtime

all-cuda: $(CUDA_VERSIONS) latest devel runtime

# cuDNN images
cudnn: cudnn-devel
	$(DOCKER_BIN) tag -f cuda:$< cuda:$@

cudnn-devel: $(CUDNN_DEVEL_LATEST)
	$(DOCKER_BIN) tag -f cuda:$< cuda:$@

cudnn-runtime: $(CUDNN_RUNTIME_LATEST)
	$(DOCKER_BIN) tag -f cuda:$< cuda:$@

# Special rules for specific cuDNN versions
%-cudnn2-devel: %-devel $(OS)/cuda/%/devel/cudnn2/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/devel/cudnn2

%-cudnn2-runtime: %-runtime $(OS)/cuda/%/runtime/cudnn2/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/runtime/cudnn2

%-cudnn3-devel: %-devel $(OS)/cuda/%/devel/cudnn3/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/devel/cudnn3

%-cudnn3-runtime: %-runtime $(OS)/cuda/%/runtime/cudnn3/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/runtime/cudnn3

%-cudnn4-devel: %-devel $(OS)/cuda/%/devel/cudnn4/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/devel/cudnn4

%-cudnn4-runtime: %-runtime $(OS)/cuda/%/runtime/cudnn4/Dockerfile
	$(DOCKER_BIN) build -t cuda:$@ $(OS)/cuda/$*/runtime/cudnn4

all-cudnn: $(CUDNN_VERSIONS) cudnn cudnn-devel cudnn-runtime

# caffe-nv images
caffe: $(OS)/caffe/Dockerfile
	$(DOCKER_BIN) build -t caffe $(OS)/caffe

push: all-cuda all-cudnn
	if [ -z "$(USERNAME)" ]; then \
		echo "Error: USERNAME not set"; \
		exit 1; \
	fi; \
        # Retag all images with the username as a prefix.
	$(DOCKER_BIN) images | awk '$$1 == "cuda" { print $$2 }' | xargs -I {} $(DOCKER_BIN) tag -f cuda:{} $(USERNAME)/cuda:{}
	$(DOCKER_BIN) push $(USERNAME)/cuda
	$(DOCKER_BIN) images | awk '$$1 == "$(USERNAME)/cuda" { print $$2 }' | xargs -I {} $(DOCKER_BIN) rmi $(USERNAME)/cuda:{}

pull:
	if [ -z "$(USERNAME)" ]; then \
		echo "Error: USERNAME not set"; \
		exit 1; \
	fi; \
        # Download all images from the Docker Hub and retag them to remove the prefix.
	$(DOCKER_BIN) pull --all-tags $(USERNAME)/cuda
	$(DOCKER_BIN) images | awk '$$1 == "$(USERNAME)/cuda" { print $$2 }' | \
		xargs -I {} sh -c '$(DOCKER_BIN) tag -f $(USERNAME)/cuda:{} cuda:{} ; $(DOCKER_BIN) rmi $(USERNAME)/cuda:{}'

clean:
	$(DOCKER_BIN) rmi -f `$(DOCKER_BIN) images -q --filter "label=com.nvidia.cuda.version"`
