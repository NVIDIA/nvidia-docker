# Copyright (c) 2017, NVIDIA CORPORATION. All rights reserved.

DOCKER ?= docker

VERSION := 2.0.2
PKG_REV := 1
RUNTIME_VERSION := 1.1.1

DIST_DIR  := $(CURDIR)/dist

.NOTPARALLEL:
.PHONY: all

all: xenial centos7 stretch

xenial: 17.12.0-xenial 17.09.1-xenial 17.09.0-xenial 17.06.2-xenial 17.03.2-xenial 1.13.1-xenial 1.12.6-xenial

centos7: 17.12.0.ce-centos7 17.09.1.ce-centos7 17.09.0.ce-centos7 17.06.2.ce-centos7 17.03.2.ce-centos7 1.13.1-centos7 1.12.6-centos7

stretch: 17.12.0-stretch 17.09.1-stretch 17.09.0-stretch 17.06.2-stretch 17.03.2-stretch

17.12.0-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.12.0-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.12.0~ce-0~ubuntu) | docker-ee (= 17.12.0~ee-0~ubuntu)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.12.0" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.09.1-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.09.1-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.09.1~ce-0~ubuntu) | docker-ee (= 17.09.1~ee-0~ubuntu)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.09.1" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.09.0-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.09.0-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.09.0~ce-0~ubuntu) | docker-ee (= 17.09.0~ee-0~ubuntu)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.09.0" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.06.2-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.06.2-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.06.2~ce-0~ubuntu) | docker-ee (= 17.06.2~ee-0~ubuntu)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.06.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.03.2-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.03.2-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.03.2~ce-0~ubuntu-xenial) | docker-ee (= 17.03.2~ee-0~ubuntu-xenial)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.03.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

1.13.1-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker1.13.1-1" \
                        --build-arg DOCKER_VERSION="docker-engine (= 1.13.1-0~ubuntu-xenial) | docker.io (= 1.13.1-0ubuntu1~16.04.2)" \
                        --build-arg PKG_VERS="$(VERSION)+docker1.13.1" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

1.12.6-xenial:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker1.12.6-1" \
                        --build-arg DOCKER_VERSION="docker-engine (= 1.12.6-0~ubuntu-xenial) | docker.io (= 1.12.6-0ubuntu1~16.04.1)" \
                        --build-arg PKG_VERS="$(VERSION)+docker1.12.6" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.xenial .
	$(DOCKER) run --rm -v $(DIST_DIR)/xenial:/dist:Z nvidia-docker2:$@

17.12.0.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.12.0" \
                        --build-arg DOCKER_VERSION="docker-ce = 17.12.0.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.12.0.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

17.09.1.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.09.1" \
                        --build-arg DOCKER_VERSION="docker-ce = 17.09.1.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.09.1.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

17.09.0.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.09.0" \
                        --build-arg DOCKER_VERSION="docker-ce = 17.09.0.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.09.0.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

17.06.2.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.06.2" \
                        --build-arg DOCKER_VERSION="docker-ce = 17.06.2.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.06.2.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

17.03.2.ce-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker17.03.2" \
                        --build-arg DOCKER_VERSION="docker-ce = 17.03.2.ce" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker17.03.2.ce" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

1.13.1-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker1.13.1" \
                        --build-arg DOCKER_VERSION="docker = 2:1.13.1" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker1.13.1" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

1.12.6-centos7:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)-1.docker1.12.6" \
                        --build-arg DOCKER_VERSION="docker = 2:1.12.6" \
                        --build-arg PKG_VERS="$(VERSION)" \
                        --build-arg PKG_REV="$(PKG_REV).docker1.12.6" \
                        -t nvidia-docker2:$@ -f Dockerfile.centos7 .
	$(DOCKER) run --rm -v $(DIST_DIR)/centos7:/dist:Z nvidia-docker2:$@

17.12.0-stretch:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.12.0-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.12.0~ce-0~debian)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.12.0" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.stretch .
	$(DOCKER) run --rm -v $(DIST_DIR)/stretch:/dist:Z nvidia-docker2:$@

17.09.1-stretch:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.09.1-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.09.1~ce-0~debian)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.09.1" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.stretch .
	$(DOCKER) run --rm -v $(DIST_DIR)/stretch:/dist:Z nvidia-docker2:$@

17.09.0-stretch:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.09.0-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.09.0~ce-0~debian)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.09.0" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.stretch .
	$(DOCKER) run --rm -v $(DIST_DIR)/stretch:/dist:Z nvidia-docker2:$@

17.06.2-stretch:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.06.2-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.06.2~ce-0~debian)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.06.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.stretch .
	$(DOCKER) run --rm -v $(DIST_DIR)/stretch:/dist:Z nvidia-docker2:$@

17.03.2-stretch:
	$(DOCKER) build --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)+docker17.03.2-1" \
                        --build-arg DOCKER_VERSION="docker-ce (= 17.03.2~ce-0~debian-stretch)" \
                        --build-arg PKG_VERS="$(VERSION)+docker17.03.2" \
                        --build-arg PKG_REV="$(PKG_REV)" \
                        -t nvidia-docker2:$@ -f Dockerfile.stretch .
	$(DOCKER) run --rm -v $(DIST_DIR)/stretch:/dist:Z nvidia-docker2:$@
