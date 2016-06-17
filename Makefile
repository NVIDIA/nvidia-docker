# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

OS ?= ubuntu-14.04
USERNAME ?= nvidia
WITH_PUSH_SUFFIX ?= 0
ifeq ($(WITH_PUSH_SUFFIX), 1)
	PUSH_SUFFIX := -$(subst -,,$(OS))
endif

.NOTPARALLEL:
.PHONY: tools clean install cuda caffe digits samples clean-images push pull

tools:
	make -C $(CURDIR)/tools

clean: clean-images
	make -C $(CURDIR)/tools clean

install:
	make -C $(CURDIR)/tools install

uninstall:
	make -C $(CURDIR)/tools uninstall

deb:
	make -C $(CURDIR)/tools deb

rpm:
	make -C $(CURDIR)/tools rpm

tarball:
	make -C $(CURDIR)/tools tarball

cuda: $(CURDIR)/$(OS)/cuda
	make -C $(CURDIR)/$(OS)/cuda

opencl: $(CURDIR)/$(OS)/opencl
	make -C $(CURDIR)/$(OS)/opencl

caffe: $(CURDIR)/$(OS)/caffe
	make -C $(CURDIR)/$(OS)/caffe

digits: $(CURDIR)/$(OS)/digits
	make -C $(CURDIR)/$(OS)/digits

samples: $(CURDIR)/samples
	make -C $(CURDIR)/$(OS)/cuda latest
	make -C $(CURDIR)/samples/$(OS)

rm_images = \
$(NV_DOCKER) images | awk '$$1 == "$(1)" {print $$1":"$$2}' | xargs -r $(NV_DOCKER) rmi

clean-images:
	$(call rm_images,cuda)
	$(call rm_images,caffe)
	$(call rm_images,digits)

# Tag all images with the Docker Hub username as a prefix, push them and untag everything.
dockerhub_push = \
$(NV_DOCKER) images | awk '$$1 == "$(1)" {print $$1":"$$2}' | xargs -I{} $(NV_DOCKER) tag {} $(USERNAME)/{}$(PUSH_SUFFIX) && \
($(NV_DOCKER) push $(USERNAME)/$(1) || true) && \
$(NV_DOCKER) images | awk '$$1 == "$(USERNAME)/$(1)" {print $$1":"$$2}' | xargs -r $(NV_DOCKER) rmi

# Download all images from the Docker Hub and retag them to remove the prefix.
dockerhub_pull = \
$(NV_DOCKER) pull --all-tags $(USERNAME)/$(1) && \
$(NV_DOCKER) images | awk '$$1 == "$(USERNAME)/$(1)" {print $$2}' | \
  xargs -I{} sh -c '$(NV_DOCKER) tag $(USERNAME)/$(1):{} $(1):{} ; $(NV_DOCKER) rmi $(USERNAME)/$(1):{}'

push:
	$(call dockerhub_push,cuda)
	$(call dockerhub_push,caffe)
	$(call dockerhub_push,digits)

pull:
	$(call dockerhub_pull,cuda)
	$(call dockerhub_pull,caffe)
	$(call dockerhub_pull,digits)
