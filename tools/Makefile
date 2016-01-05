# Copyright (c) 2015, NVIDIA CORPORATION. All rights reserved.

BIN_DIR  := $(CURDIR)/bin
USER_ID  := $(shell id -u)

IMAGE  := nvdocker-build
PREFIX := /usr/bin
DOCKER := $(BIN_DIR)/nvidia-docker
PLUGIN := $(BIN_DIR)/nvidia-docker-plugin

.PHONY: all install clean

all : $(BIN_DIR)

$(BIN_DIR) :
	@docker build --build-arg UID=$(USER_ID) -t $(IMAGE) -f Dockerfile.build $(CURDIR)
	@mkdir -p $(BIN_DIR)
	@docker run --rm -v $(BIN_DIR):/go/bin $(IMAGE)

install: all
	install -D -T -m 755 $(DOCKER) $(PREFIX)/$(notdir $(DOCKER))
	install -D -T -m 755 $(PLUGIN) $(PREFIX)/$(notdir $(PLUGIN))

uninstall:
	$(RM) $(PREFIX)/$(notdir $(DOCKER)) $(PREFIX)/$(notdir $(PLUGIN))

clean :
	-@docker rmi -f $(IMAGE) golang 2> /dev/null
	@rm -rf $(BIN_DIR)
