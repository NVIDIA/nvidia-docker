# Copyright (c) 2017, NVIDIA CORPORATION. All rights reserved.

DOCKER ?= docker
MKDIR  ?= mkdir
DIST_DIR ?= $(CURDIR)/dist

LIB_NAME := nvidia-docker2
LIB_VERSION := 2.4.0
PKG_REV := 1

RUNTIME_VERSION := 3.3.0

# Supported OSs by architecture
AMD64_TARGETS := ubuntu20.04 ubuntu18.04 ubuntu16.04 debian10 debian9
X86_64_TARGETS := centos7 centos8 rhel7 rhel8 amazonlinux1 amazonlinux2 opensuse-leap15.1
PPC64LE_TARGETS := ubuntu18.04 ubuntu16.04 centos7 centos8 rhel7 rhel8
ARM64_TARGETS := ubuntu20.04 ubuntu18.04
AARCH64_TARGETS := centos8 rhel8

# By default run all native docker-based targets
docker-native:

# Define top-level build targets
docker%: SHELL:=/bin/bash

# Native targets
PLATFORM ?= $(shell uname -m)
ifeq ($(PLATFORM),x86_64)
NATIVE_TARGETS := $(AMD64_TARGETS) $(X86_64_TARGETS)
$(AMD64_TARGETS): %: %-amd64
$(X86_64_TARGETS): %: %-x86_64
else ifeq ($(PLATFORM),ppc64le)
NATIVE_TARGETS := $(PPC64LE_TARGETS)
$(PPC64LE_TARGETS): %: %-ppc64le
else ifeq ($(PLATFORM),aarch64)
NATIVE_TARGETS := $(ARM64_TARGETS) $(AARCH64_TARGETS)
$(ARM64_TARGETS): %: %-arm64
$(AARCH64_TARGETS): %: %-aarch64
endif
docker-native: $(NATIVE_TARGETS)

# amd64 targets
AMD64_TARGETS := $(patsubst %, %-amd64, $(AMD64_TARGETS))
$(AMD64_TARGETS): ARCH := amd64
$(AMD64_TARGETS): %: --%
docker-amd64: $(AMD64_TARGETS)

# x86_64 targets
X86_64_TARGETS := $(patsubst %, %-x86_64, $(X86_64_TARGETS))
$(X86_64_TARGETS): ARCH := x86_64
$(X86_64_TARGETS): %: --%
docker-x86_64: $(X86_64_TARGETS)

# arm64 targets
ARM64_TARGETS := $(patsubst %, %-arm64, $(ARM64_TARGETS))
$(ARM64_TARGETS): ARCH := arm64
$(ARM64_TARGETS): %: --%
docker-arm64: $(ARM64_TARGETS)

# aarch64 targets
AARCH64_TARGETS := $(patsubst %, %-aarch64, $(AARCH64_TARGETS))
$(AARCH64_TARGETS): ARCH := aarch64
$(AARCH64_TARGETS): %: --%
docker-aarch64: $(AARCH64_TARGETS)

# ppc64le targets
PPC64LE_TARGETS := $(patsubst %, %-ppc64le, $(PPC64LE_TARGETS))
$(PPC64LE_TARGETS): ARCH := ppc64le
$(PPC64LE_TARGETS): WITH_LIBELF := yes
$(PPC64LE_TARGETS): %: --%
docker-ppc64le: $(PPC64LE_TARGETS)

# docker target to build for all os/arch combinations
docker-all: $(AMD64_TARGETS) $(X86_64_TARGETS) \
            $(ARM64_TARGETS) $(AARCH64_TARGETS) \
            $(PPC64LE_TARGETS)

# Default variables for all private '--' targets below.
# One private target is defined for each OS we support.
--%: TARGET_PLATFORM = $(*)
--%: VERSION = $(patsubst $(OS)%-$(ARCH),%,$(TARGET_PLATFORM))
--%: BASEIMAGE = $(OS):$(VERSION)
--%: BUILDIMAGE = nvidia/$(LIB_NAME)/$(OS)$(VERSION)-$(ARCH)
--%: DOCKERFILE = $(CURDIR)/docker/Dockerfile.$(OS)
--%: ARTIFACTS_DIR = $(DIST_DIR)/$(OS)$(VERSION)/$(ARCH)
--%: docker-build-%
	@

# private ubuntu target
--ubuntu%: OS := ubuntu
--ubuntu%: DOCKER_VERSION := docker-ce (>= 18.06.0~ce~3-0~ubuntu) | docker-ee (>= 18.06.0~ce~3-0~ubuntu) | docker.io (>= 18.06.0)

# private debian target
--debian%: OS := debian
--debian%: DOCKER_VERSION := docker-ce (>= 18.06.0~ce~3-0~debian) | docker-ee (>= 18.06.0~ce~3-0~debian) | docker.io (>= 18.06.0)

# private centos target
--centos%: OS := centos
--centos%: DOCKER_VERSION := docker-ce >= 18.06.3.ce-3.el7

# private amazonlinuxtarget
--amazonlinux%: OS := amazonlinux
--amazonlinux2%: DOCKER_VERSION := docker >= 18.06.1ce-2.amzn2
--amazonlinux1%: DOCKER_VERSION := docker >= 18.06.1ce-2.16.amzn1

# private opensuse-leap target
--opensuse-leap%: OS := opensuse-leap
--opensuse-leap%: BASEIMAGE = opensuse/leap:$(VERSION)
--opensuse-leap%: DOCKER_VERSION := docker >= 18.09.1_ce

# private rhel target (actually built on centos)
--rhel%: OS := centos
--rhel%: VERSION = $(patsubst rhel%-$(ARCH),%,$(TARGET_PLATFORM))
--rhel%: ARTIFACTS_DIR = $(DIST_DIR)/rhel$(VERSION)/$(ARCH)
--rhel%: DOCKER_VERSION := docker-ce >= 18.06.3.ce-3.el7

docker-build-%:
	@echo "Building for $(TARGET_PLATFORM)"
	docker pull --platform=linux/$(ARCH) $(BASEIMAGE)
	DOCKER_BUILDKIT=1 \
	$(DOCKER) build \
	    --progress=plain \
	    --build-arg BASEIMAGE=$(BASEIMAGE) \
	    --build-arg DOCKER_VERSION="$(DOCKER_VERSION)" \
	    --build-arg RUNTIME_VERSION="$(RUNTIME_VERSION)" \
	    --build-arg PKG_VERS="$(LIB_VERSION)" \
	    --build-arg PKG_REV="$(PKG_REV)" \
	    --tag $(BUILDIMAGE) \
	    --file $(DOCKERFILE) .
	$(DOCKER) run \
	    -e DISTRIB \
	    -e SECTION \
	    -v $(ARTIFACTS_DIR):/dist \
	    $(BUILDIMAGE)

docker-clean:
	IMAGES=$$(docker images "nvidia/$(LIB_NAME)/*" --format="{{.ID}}"); \
	if [ "$${IMAGES}" != "" ]; then \
	    docker rmi -f $${IMAGES}; \
	fi

distclean:
	rm -rf $(DIST_DIR)
