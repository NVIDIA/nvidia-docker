# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DOCKER ?= docker
MKDIR  ?= mkdir
DIST_DIR ?= $(CURDIR)/dist

# Supported OSs by architecture
AMD64_TARGETS := ubuntu20.04 ubuntu18.04 ubuntu16.04 debian10 debian9
X86_64_TARGETS := centos7 centos8 rhel7 rhel8 amazonlinux2 opensuse-leap15.1
PPC64LE_TARGETS := ubuntu18.04 ubuntu16.04 centos7 centos8 rhel7 rhel8
ARM64_TARGETS := ubuntu20.04 ubuntu18.04
AARCH64_TARGETS := centos8 rhel8 amazonlinux2

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

DEB_LIB_VERSION = $(LIB_VERSION)$(if $(LIB_TAG),~$(LIB_TAG))
DEB_PKG_REV = 1
DEB_TOOLKIT_VERSION = $(TOOLKIT_VERSION)$(if $(TOOLKIT_TAG),~$(TOOLKIT_TAG))
DEB_TOOLKIT_REV = 1

RPM_LIB_VERSION = $(LIB_VERSION)
RPM_PKG_REV = $(if $(LIB_TAG),0.1.$(LIB_TAG),1)
RPM_TOOLKIT_VERSION = $(TOOLKIT_VERSION)
RPM_TOOLKIT_REV = $(if $(TOOLKIT_TAG),0.1.$(TOOLKIT_TAG),1)

# private OS targets with defaults
# private ubuntu target
--ubuntu%: OS := ubuntu
--ubuntu%: PKG_VERS = $(DEB_LIB_VERSION)
--ubuntu%: PKG_REV = $(DEB_PKG_REV)
--ubuntu%: MIN_TOOLKIT_PKG_VERSION = $(DEB_TOOLKIT_VERSION)-$(DEB_TOOLKIT_REV)

# private debian target
--debian%: OS := debian
--debian%: PKG_VERS = $(DEB_LIB_VERSION)
--debian%: PKG_REV = $(DEB_PKG_REV)
--debian%: MIN_TOOLKIT_PKG_VERSION = $(DEB_TOOLKIT_VERSION)-$(DEB_TOOLKIT_REV)


# private centos target
--centos%: OS := centos
--centos%: PKG_VERS = $(LIB_VERSION)
--centos%: PKG_REV = $(RPM_PKG_REV)
--centos%: MIN_TOOLKIT_PKG_VERSION = $(RPM_TOOLKIT_VERSION)-$(RPM_TOOLKIT_REV)
--centos8%: BASEIMAGE = quay.io/centos/centos:stream8

# private amazonlinux target
--amazonlinux%: OS := amazonlinux
--amazonlinux%: PKG_VERS = $(LIB_VERSION)
--amazonlinux%: PKG_REV = $(RPM_PKG_REV)
--amazonlinux%: MIN_TOOLKIT_PKG_VERSION = $(RPM_TOOLKIT_VERSION)-$(RPM_TOOLKIT_REV)

# private opensuse-leap target with overrides
--opensuse-leap%: OS := opensuse-leap
--opensuse-leap%: PKG_VERS = $(LIB_VERSION)
--opensuse-leap%: PKG_REV = $(RPM_PKG_REV)
--opensuse-leap%: MIN_TOOLKIT_PKG_VERSION = $(RPM_TOOLKIT_VERSION)-$(RPM_TOOLKIT_REV)
--opensuse-leap%: BASEIMAGE = opensuse/leap:$(VERSION)

# private rhel target (actually built on centos)
--rhel%: OS := centos
--rhel%: PKG_VERS = $(LIB_VERSION)
--rhel%: PKG_REV = $(RPM_PKG_REV)
--rhel%: MIN_TOOLKIT_PKG_VERSION = $(RPM_TOOLKIT_VERSION)-$(RPM_TOOLKIT_REV)
--rhel%: VERSION = $(patsubst rhel%-$(ARCH),%,$(TARGET_PLATFORM))
--rhel%: ARTIFACTS_DIR = $(DIST_DIR)/rhel$(VERSION)/$(ARCH)
--rhel8%: BASEIMAGE = quay.io/centos/centos:stream8

# Specify required docker versions
--ubuntu%: DOCKER_VERSION := docker-ce (>= 18.06.0~ce~3-0~ubuntu) | docker-ee (>= 18.06.0~ce~3-0~ubuntu) | docker.io (>= 18.06.0) | moby-engine
--debian%: DOCKER_VERSION := docker-ce (>= 18.06.0~ce~3-0~debian) | docker-ee (>= 18.06.0~ce~3-0~debian) | docker.io (>= 18.06.0) | moby-engine
--centos%: DOCKER_VERSION := docker-ce >= 18.06.3.ce-3.el7
--amazonlinux2%: DOCKER_VERSION := docker >= 18.06.1ce-2.amzn2
--opensuse-leap%: DOCKER_VERSION := docker >= 18.09.1_ce
--rhel%: DOCKER_VERSION := docker-ce >= 18.06.3.ce-3.el7

# Depending on the docker version we may have to add the platform args to the
# build and run commands
PLATFORM_ARGS ?= --platform=linux/$(ARCH)
ifneq ($(strip $(ADD_DOCKER_PLATFORM_ARGS)),)
DOCKER_PLATFORM_ARGS = $(PLATFORM_ARGS)
endif

docker-build-%:
	@echo "Building for $(TARGET_PLATFORM)"
	docker pull $(PLATFORM_ARGS) $(BASEIMAGE)
	DOCKER_BUILDKIT=1 \
	$(DOCKER) build $(DOCKER_PLATFORM_ARGS) \
	    --progress=plain \
	    --build-arg BASEIMAGE="$(BASEIMAGE)" \
	    --build-arg DOCKER_VERSION="$(DOCKER_VERSION)" \
	    --build-arg TOOLKIT_VERSION="$(MIN_TOOLKIT_PKG_VERSION)" \
		--build-arg PKG_NAME="$(LIB_NAME)" \
	    --build-arg PKG_VERS="$(PKG_VERS)" \
	    --build-arg PKG_REV="$(PKG_REV)" \
	    --tag $(BUILDIMAGE) \
	    --file $(DOCKERFILE) .
	$(DOCKER) run $(DOCKER_PLATFORM_ARGS) \
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
