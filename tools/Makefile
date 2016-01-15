# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker
PREFIX    ?= /usr/bin

BIN_DIR  := $(CURDIR)/bin
CONF_DIR := $(CURDIR)/conf
DIST_DIR := $(CURDIR)/dist

BUILD_IMAGE := nvdocker-build
USER_ID     := $(shell id -u)

DOCKER_BIN := $(BIN_DIR)/nvidia-docker
PLUGIN_BIN := $(BIN_DIR)/nvidia-docker-plugin

PKG_VERSION := nvidia-docker_1.0.0~beta
PKG_DIR     := $(DIST_DIR)/$(subst _,-,$(PKG_VERSION))
PKG_ARCH    := amd64

.PHONY: all build install uninstall clean distclean tarball deb

all: build

build:
	@$(NV_DOCKER) build --build-arg UID=$(USER_ID) -t $(BUILD_IMAGE) -f Dockerfile.build $(CURDIR)
	@mkdir -p $(BIN_DIR)
	@$(NV_DOCKER) run --rm -v $(BIN_DIR):/go/bin $(BUILD_IMAGE)

install: all
	install -D -T -m 755 $(DOCKER_BIN) $(PREFIX)/$(notdir $(DOCKER_BIN))
	install -D -T -m 755 $(PLUGIN_BIN) $(PREFIX)/$(notdir $(PLUGIN_BIN))

uninstall:
	$(RM) $(PREFIX)/$(notdir $(DOCKER_BIN))
	$(RM) $(PREFIX)/$(notdir $(PLUGIN_BIN))

clean:
	-@$(NV_DOCKER) rmi -f $(BUILD_IMAGE) golang 2> /dev/null
	@rm -rf $(BIN_DIR)

distclean:
	@rm -rf $(DIST_DIR)

tarball: build distclean
	mkdir -p $(DIST_DIR)
	tar --transform='s;.*/;nvidia-docker/;' -caf $(DIST_DIR)/$(PKG_VERSION)_$(PKG_ARCH).tar.xz $(BIN_DIR)/*

deb: export DEBFULLNAME=NVIDIA CORPORATION
deb: export DEBEMAIL=digits@nvidia.com

deb: build distclean
	mkdir -p $(PKG_DIR)
	cp -r $(BIN_DIR) $(PKG_DIR)
	cd $(PKG_DIR) && dh_make -y -s -c bsd -p $(PKG_VERSION) -d -t $(CONF_DIR)/debian --createorig
	cd $(PKG_DIR) && debuild -e PREFIX=$(PREFIX) -i -us -uc -b
