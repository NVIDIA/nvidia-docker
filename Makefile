# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

prefix      ?= /usr/local
exec_prefix ?= $(prefix)
bindir      ?= $(exec_prefix)/bin

CR_NAME  := NVIDIA CORPORATION
CR_EMAIL := digits@nvidia.com
PKG_NAME := nvidia-docker
PKG_VERS := 1.0.1
PKG_REV  := 1
ifneq ($(MAKECMDGOALS),rpm)
PKG_ARCH := amd64
else
PKG_ARCH := x86_64
endif

# Mirror the BUILD_ARCH from the build Dockerfile
BUILD_ARCH = .$(shell uname -m)
ifneq ($(BUILD_ARCH),.ppc64le)
    BUILD_ARCH =
else
    PKG_ARCH = ppc64le
endif

BIN_DIR    := $(CURDIR)/bin
DIST_DIR   := $(CURDIR)/dist
BUILD_DIR  := $(CURDIR)/build
DOCKER_BIN := $(BIN_DIR)/nvidia-docker
PLUGIN_BIN := $(BIN_DIR)/nvidia-docker-plugin

DOCKER_VERS      := $(shell $(NV_DOCKER) version -f '{{.Client.Version}}')
DOCKER_VERS_MAJ  := $(shell echo $(DOCKER_VERS) | cut -d. -f1)
DOCKER_VERS_MIN  := $(shell echo $(DOCKER_VERS) | cut -d. -f2)

DOCKER_RMI       := $(NV_DOCKER) rmi
DOCKER_RUN       := $(NV_DOCKER) run --rm --net=host
DOCKER_IMAGES    := $(NV_DOCKER) images -q $(PKG_NAME)
DOCKER_BUILD     := $(NV_DOCKER) build --build-arg USER_ID="$(shell id -u)" \
                                       --build-arg CR_NAME="$(CR_NAME)" \
                                       --build-arg CR_EMAIL="$(CR_EMAIL)" \
                                       --build-arg PKG_NAME="$(PKG_NAME)" \
                                       --build-arg PKG_VERS="$(PKG_VERS)" \
                                       --build-arg PKG_REV="$(PKG_REV)" \
                                       --build-arg PKG_ARCH="$(PKG_ARCH)"

.PHONY: all build install uninstall clean distclean tarball deb rpm

all: build

build: distclean
	@mkdir -p $(BIN_DIR)
	@$(DOCKER_BUILD) -t $(PKG_NAME):$@ -f Dockerfile.$@$(BUILD_ARCH) $(CURDIR)
	@$(DOCKER_RUN) -v $(BIN_DIR):/go/bin:Z $(PKG_NAME):$@

install: build
	install -D -m 755 -t $(bindir) $(DOCKER_BIN)
	install -D -m 755 -t $(bindir) $(PLUGIN_BIN)

uninstall:
	$(RM) $(bindir)/$(notdir $(DOCKER_BIN))
	$(RM) $(bindir)/$(notdir $(PLUGIN_BIN))

clean:
	-@$(DOCKER_IMAGES) | xargs $(DOCKER_RMI) 2> /dev/null
	-@$(DOCKER_RMI) golang:1.5 ubuntu:14.04 centos:7 2> /dev/null

distclean:
	@rm -rf $(BIN_DIR)
	@rm -rf $(DIST_DIR)

tarball: build
	@mkdir -p $(DIST_DIR)
	tar --transform='s;.*/;$(PKG_NAME)/;' -caf $(DIST_DIR)/$(PKG_NAME)_$(PKG_VERS)_$(PKG_ARCH).tar.xz $(BIN_DIR)/*
	@printf "\nFind tarball at $(DIST_DIR)\n\n"

deb: tarball
	@$(DOCKER_BUILD) -t $(PKG_NAME):$@ -f Dockerfile.$@$(BUILD_ARCH) $(CURDIR)
	@$(DOCKER_RUN) -ti -v $(DIST_DIR):/dist:Z -v $(BUILD_DIR):/build:Z $(PKG_NAME):$@
	@printf "\nFind packages at $(DIST_DIR)\n\n"

rpm: tarball
	@$(DOCKER_BUILD) -t $(PKG_NAME):$@ -f Dockerfile.$@$(BUILD_ARCH) $(CURDIR)
	@$(DOCKER_RUN) -ti -v $(DIST_DIR):/dist:Z -v $(BUILD_DIR):/build:Z $(PKG_NAME):$@
	@printf "\nFind packages at $(DIST_DIR)\n\n"
