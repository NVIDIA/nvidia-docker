# Copyright (c) 2017, NVIDIA CORPORATION. All rights reserved.

DOCKER ?= docker

VERSION := 2.0.0
PKG_REV := 1
RUNTIME_VERSION := 1.0.0

DIST_DIR  := $(CURDIR)/dist

.NOTPARALLEL:
.PHONY: all

all: xenial centos7

xenial: 17.06.2-xenial 17.06.1-xenial 17.03.2-xenial 1.13.1-xenial 1.12.6-xenial

centos7: 17.06.2.ce-centos7 17.06.1.ce-centos7 17.03.2.ce-centos7

17.06.2-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.06.2-1" \
                        --build-arg DOCKER_VERSION="17.06.2~ce-0~ubuntu" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.06.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.06.1-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.06.1-1" \
                        --build-arg DOCKER_VERSION="17.06.1~ce-0~ubuntu" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.06.1" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.03.2-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.03.2-1" \
                        --build-arg DOCKER_VERSION="17.03.2~ce-0~ubuntu-xenial" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.03.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

1.13.1-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker1.13.1-1" \
                        --build-arg DOCKER_VERSION="1.13.1-0~ubuntu-xenial" \
                        --build-arg PKG_VERS="$(VERSION)+docker1.13.1" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

1.12.6-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker1.12.6-1" \
                        --build-arg DOCKER_VERSION="1.12.6-0~ubuntu-xenial" \
                        --build-arg PKG_VERS="$(VERSION)+docker1.12.6" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.06.2.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.06.2" \
                        --build-arg DOCKER_VERSION="17.06.2.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.06.2.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

17.06.1.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.06.1" \
                        --build-arg DOCKER_VERSION="17.06.1.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.06.1.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

17.03.2.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.03.2" \
                        --build-arg DOCKER_VERSION="17.03.2.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.03.2.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@
