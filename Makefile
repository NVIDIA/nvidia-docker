# Copyright (c) 2017, NVIDIA CORPORATION. All rights reserved.

DOCKER ?= docker

VERSION := 2.0.3
PKG_REV := 1
RUNTIME_VERSION := 2.0.0

DIST_DIR  := $(CURDIR)/dist

.NOTPARALLEL:
.PHONY: all

all: ubuntu16.04 debian9 centos7 amzn1

ubuntu16.04: $(addsuffix -ubuntu16.04, 17.12.1 17.12.0 17.09.1 17.09.0 17.06.2 17.03.2 1.13.1 1.12.6)

debian9: $(addsuffix -debian9, 17.12.1 17.12.0 17.09.1 17.09.0 17.06.2 17.03.2)

centos7: $(addsuffix -centos7, 17.12.1.ce 17.12.0.ce 17.09.1.ce 17.09.0.ce 17.06.2.ce 17.03.2.ce 1.13.1 1.12.6)

amzn1: $(addsuffix -amzn1, 17.09.1.ce 17.06.2.ce 17.03.2.ce)


%-ubuntu16.04:
	$(DOCKER) build --build-arg VERSION_ID="16.04" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker$*-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= $*~ce-0~ubuntu) | docker-ee (= $*~ee-0~ubuntu)" \
                        --build-arg PKG_VERS="$(VERSION)+docker$*" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/ubuntu:16.04-docker$*" -f Dockerfile.ubuntu .
	$(DOCKER) run --rm -v $(DIST_DIR)/ubuntu16.04:/dist:Z "nvidia/nvidia-docker2/ubuntu:16.04-docker$*"

17.03.2-ubuntu16.04:
	$(DOCKER) build --build-arg VERSION_ID="16.04" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.03.2-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.03.2~ce-0~ubuntu-xenial) | docker-ee (= 17.03.2~ee-0~ubuntu-xenial)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.03.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/ubuntu:16.04-docker$*" -f Dockerfile.ubuntu .
	$(DOCKER) run --rm -v $(DIST_DIR)/ubuntu16.04:/dist:Z "nvidia/nvidia-docker2/ubuntu:16.04-docker$*"

1.13.1-ubuntu16.04:
	$(DOCKER) build --build-arg VERSION_ID="16.04" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker1.13.1-1" \
                        --build-arg DOCKER_VERSION="docker-engine (= 1.13.1-0~ubuntu-xenial) | docker.io (= 1.13.1-0ubuntu1~16.04.2)" \
                        --build-arg PKG_VERS="$(VERSION)+docker1.13.1" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/ubuntu:16.04-docker1.13.1" -f Dockerfile.ubuntu .
	$(DOCKER) run --rm -v $(DIST_DIR)/ubuntu16.04:/dist:Z "nvidia/nvidia-docker2/ubuntu:16.04-docker1.13.1"

1.12.6-ubuntu16.04:
	$(DOCKER) build --build-arg VERSION_ID="16.04" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker1.12.6-1" \
                        --build-arg DOCKER_VERSION="docker-engine (= 1.12.6-0~ubuntu-xenial) | docker.io (= 1.12.6-0ubuntu1~16.04.1)" \
                        --build-arg PKG_VERS="$(VERSION)+docker1.12.6" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/ubuntu:16.04-docker1.12.6" -f Dockerfile.ubuntu .
	$(DOCKER) run --rm -v $(DIST_DIR)/ubuntu16.04:/dist:Z "nvidia/nvidia-docker2/ubuntu:16.04-docker1.12.6"

%-debian9:
	$(DOCKER) build --build-arg VERSION_ID="9" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker$*-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= $*~ce-0~debian)" \
                        --build-arg PKG_VERS="$(VERSION)+docker$*" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/debian:9-docker$*" -f Dockerfile.debian .
	$(DOCKER) run --rm -v $(DIST_DIR)/debian9:/dist:Z "nvidia/nvidia-docker2/debian:9-docker$*"

17.03.2-debian9:
	$(DOCKER) build --build-arg VERSION_ID="9" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.03.2-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.03.2~ce-0~debian-stretch)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.03.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t "nvidia/nvidia-docker2/debian:9-docker17.03.2" -f Dockerfile.debian .
	$(DOCKER) run --rm -v $(DIST_DIR)/debian9:/dist:Z "nvidia/nvidia-docker2/debian:9-docker17.03.2"

%.ce-centos7:
	$(DOCKER) build --build-arg VERSION_ID="7" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker$*" \
                        --build-arg DOCKER_VERSION="docker-ce = $*.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker$*.ce" \
                        -t "nvidia/nvidia-docker2/centos:7-docker$*.ce" -f Dockerfile.centos .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z "nvidia/nvidia-docker2/centos:7-docker$*.ce"

%-centos7:
	$(DOCKER) build --build-arg VERSION_ID="7" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker$*" \
                        --build-arg DOCKER_VERSION="docker = 2:$*" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker$*" \
                        -t "nvidia/nvidia-docker2/centos:7-docker$*" -f Dockerfile.centos .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z "nvidia/nvidia-docker2/centos:7-docker$*"

%.ce-amzn1:
	$(DOCKER) build --build-arg VERSION_ID="1" \
                        --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker$*.amzn1" \
                        --build-arg DOCKER_VERSION="docker = $*ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker$*.ce.amzn1" \
                        -t "nvidia/nvidia-docker2/amzn:1-docker$*" -f Dockerfile.amzn .
	$(DOCKER) run --rm -v $(DIST_DIR)/amzn1:/dist:Z "nvidia/nvidia-docker2/amzn:1-docker$*"
