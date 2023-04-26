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

# Supported packaging formats
FORMAT_TARGETS := deb rpm

# We add utility targets to support common os-arch combinations by mapping to the required format targets.
DEB_TARGETS := debian10-amd64 ubuntu18.04-amd64 ubuntu18.04-arm64 ubuntu18.04-ppc64le
RPM_TARGETS := amazonlinux2-aarch64 amazonlinux2-x86_64 centos7-x86_64 centos8-aarch64 centos8-ppc64le centos8-x86_64 opensuse-leap15.1-x86_64

$(DEB_TARGETS): %: deb
$(RPM_TARGETS): %: rpm

# Define top-level build targets
docker%: SHELL:=/bin/bash

$(FORMAT_TARGETS): %: --%

# Default variables for all private '--' targets below.
# One private target is defined for each OS we support.
--%: FORMAT = $(*)
--%: BUILDIMAGE = nvidia/$(LIB_NAME)/$(FORMAT)-all
--%: DOCKERFILE = $(CURDIR)/docker/Dockerfile.$(FORMAT)
--%: ARTIFACTS_DIR = $(DIST_DIR)/$(FORMAT)/all
--%: docker-build-%
	@

PKG_VERS = $(LIB_VERSION)$(if $(LIB_TAG),~$(LIB_TAG))
PKG_REV = 1
MIN_TOOLKIT_PKG_VERSION = $(TOOLKIT_VERSION)$(if $(TOOLKIT_TAG),~$(TOOLKIT_TAG))-1

--deb: BASEIMAGE := ubuntu:18.04

--rpm: BASEIMAGE := quay.io/centos/centos:stream8

docker-build-%:
	@echo "Building $(FORMAT) packages to $(ARTIFACTS_DIR)"
	docker pull $(BASEIMAGE)
	DOCKER_BUILDKIT=1 \
	$(DOCKER) build \
	    --progress=plain \
	    --build-arg BASEIMAGE="$(BASEIMAGE)" \
	    --build-arg TOOLKIT_VERSION="$(MIN_TOOLKIT_PKG_VERSION)" \
	    --build-arg PKG_NAME="$(LIB_NAME)" \
	    --build-arg PKG_VERS="$(PKG_VERS)" \
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
