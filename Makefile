# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

OS ?= ubuntu
USERNAME ?= nvidia

.PHONY: tools clean install cuda caffe digits samples push pull

tools:
	make -C $(CURDIR)/tools

clean:
	make -C $(CURDIR)/tools clean

install:
	make -C $(CURDIR)/tools install

uninstall:
	make -C $(CURDIR)/tools uninstall

cuda: $(CURDIR)/$(OS)/cuda
	make -C $(CURDIR)/$(OS)/cuda

caffe: $(CURDIR)/$(OS)/caffe
	make -C $(CURDIR)/$(OS)/caffe

digits: $(CURDIR)/$(OS)/digits
	make -C $(CURDIR)/$(OS)/digits

samples: $(CURDIR)/samples
	make -C $(CURDIR)/$(OS)/cuda latest
	make -C $(CURDIR)/samples/$(OS)

# Tag all images with the Docker Hub username as a prefix, push them and untag everything.
dockerhub_push = \
$(NV_DOCKER) images | awk '$$1 == "$(1)" {print $$1":"$$2}' | xargs -I{} $(NV_DOCKER) tag -f {} $(USERNAME)/{} && \
$(NV_DOCKER) push $(USERNAME)/$(1) && \
$(NV_DOCKER) images | awk '$$1 == "$(USERNAME)/$(1)" {print $$1":"$$2}' | xargs $(NV_DOCKER) rmi

# Download all images from the Docker Hub and retag them to remove the prefix.
dockerhub_pull = \
$(NV_DOCKER) pull --all-tags $(USERNAME)/$(1) && \
$(NV_DOCKER) images | awk '$$1 == "$(USERNAME)/$(1)" {print $$2}' | \
  xargs -I{} sh -c '$(NV_DOCKER) tag -f $(USERNAME)/$(1):{} $(1):{} ; $(NV_DOCKER) rmi $(USERNAME)/$(1):{}'

push: cuda caffe digits
	$(call dockerhub_push,cuda)
	$(call dockerhub_push,caffe)
	$(call dockerhub_push,digits)

pull:
	$(call dockerhub_pull,cuda)
	$(call dockerhub_pull,caffe)
	$(call dockerhub_pull,digits)
