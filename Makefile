# Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

NV_DOCKER ?= docker

.NOTPARALLEL:
.PHONY: tools clean install uninstall deb rpm tarball samples

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

samples:
	make -C $(CURDIR)/samples/ubuntu-16.04
