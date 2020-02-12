# Copyright (c) 2017, NVIDIA CORPORATION. All rights reserved.

DOCKER ?= docker
MKDIR  ?= mkdir

VERSION := 2.2.2
RUNTIME_VERSION := 3.1.2
PKG_REV := 1

DIST_DIR  := $(CURDIR)/dist

.NOTPARALLEL:
.PHONY: all

all: ubuntu18.04-amd64 ubuntu16.04-amd64 debian10-amd64 debian9-amd64 centos7-x86_64 amzn2-x86_64 amzn1-x86_64 opensuse-leap15.1-x86_64 \
     ubuntu16.04-ppc64le ubuntu18.04-ppc64le centos7-ppc64le

ubuntu18.04-%:
	$(DOCKER) build --build-arg VERSION_ID="18.04" \
                        --build-arg DOCKER_VERSION="docker-ce (>= 18.06.0~ce~3-0~ubuntu) | docker-ee (>= 18.06.0~ce~3-0~ubuntu) | docker.io (>= 18.06.0)" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/ubuntu:18.04-$*" -f docker/$*/Dockerfile.ubuntu .
	$(MKDIR) -p $(DIST_DIR)/ubuntu18.04/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/ubuntu:18.04-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/ubuntu18.04/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid

ubuntu16.04-%:
	$(DOCKER) build --build-arg VERSION_ID="16.04" \
                        --build-arg DOCKER_VERSION="docker-ce (>= 18.06.0~ce~3-0~ubuntu) | docker-ee (>= 18.06.0~ce~3-0~ubuntu) | docker.io (>= 18.06.0)" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/ubuntu:16.04-$*" -f docker/$*/Dockerfile.ubuntu .
	$(MKDIR) -p $(DIST_DIR)/ubuntu16.04/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/ubuntu:16.04-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/ubuntu16.04/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid

debian10-%:
	$(DOCKER) build --build-arg VERSION_ID="10" \
                        --build-arg DOCKER_VERSION="docker-ce (>= 18.06.0~ce~3-0~debian) | docker-ee (>= 18.06.0~ce~3-0~debian) | docker.io (>= 18.06.0)" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/debian:10-$*" -f docker/$*/Dockerfile.debian .
	$(MKDIR) -p $(DIST_DIR)/debian10/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/debian:10-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/debian10/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid

debian9-%:
	$(DOCKER) build --build-arg VERSION_ID="9" \
                        --build-arg DOCKER_VERSION="docker-ce (>= 18.06.0~ce~3-0~debian) | docker-ee (>= 18.06.0~ce~3-0~debian) | docker.io (>= 18.06.0)" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/debian:9-$*" -f docker/$*/Dockerfile.debian .
	$(MKDIR) -p $(DIST_DIR)/debian9/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/debian:9-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/debian9/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid

centos7-%:
	$(DOCKER) build --build-arg VERSION_ID="7" \
                        --build-arg DOCKER_VERSION="docker-ce >= 18.06.3.ce-3.el7" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/centos:7-$*" -f docker/$*/Dockerfile.centos .
	$(MKDIR) -p $(DIST_DIR)/centos7/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/centos:7-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/centos7/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid

amzn2-%:
	$(DOCKER) build --build-arg VERSION_ID="2" \
                        --build-arg DOCKER_VERSION="docker >= 18.06.1ce-2.amzn2" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/amzn:2-docker-$*" -f docker/$*/Dockerfile.amzn .
	$(MKDIR) -p $(DIST_DIR)/amzn2/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/amzn:2-docker-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/amzn2/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid

amzn1-%:
	$(DOCKER) build --build-arg VERSION_ID="1" \
                        --build-arg DOCKER_VERSION="docker >= 18.06.1ce-2.16.amzn1" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/amzn:1-$*" -f docker/$*/Dockerfile.amzn .
	$(MKDIR) -p $(DIST_DIR)/amzn1/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/amzn:1-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/amzn1/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid

opensuse-leap15.1-%:
	$(DOCKER) build --build-arg VERSION_ID="15.1" \
                        --build-arg DOCKER_VERSION="docker >= 18.09.1_ce" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/opensuse-leap:15.1-$*" -f docker/$*/Dockerfile.opensuse-leap .
	$(MKDIR) -p $(DIST_DIR)/opensuse-leap15.1/$*
	$(DOCKER) run  --cidfile $@.cid "nvidia/nvidia-docker2/opensuse-leap:15.1-$*"
	$(DOCKER) cp $$(cat $@.cid):/dist/. $(DIST_DIR)/opensuse-leap15.1/$*
	$(DOCKER) rm $$(cat $@.cid) && rm $@.cid
