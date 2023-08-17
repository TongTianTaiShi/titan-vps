SHELL=/usr/bin/env bash

all: build
.PHONY: all

unexport GOFLAGS

GOCC?=go

GOVERSION:=$(shell $(GOCC) version | tr ' ' '\n' | grep go1 | sed 's/^go//' | awk -F. '{printf "%d%03d%03d", $$1, $$2, $$3}')
ifeq ($(shell expr $(GOVERSION) \< 1017001), 1)
$(warning Your Golang version is go$(shell expr $(GOVERSION) / 1000000).$(shell expr $(GOVERSION) % 1000000 / 1000).$(shell expr $(GOVERSION) % 1000))
$(error Update Golang to version to at least 1.18)
endif

# git modules that need to be loaded
MODULES:=

CLEAN:=
BINS:=

ldflags=-X=github.com/LMF709268224/titan-vps/build.CurrentCommit=+git.$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))
ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"


mall: $(BUILD_DEPS)
	rm -f mall
	$(GOCC) build $(GOFLAGS) -o mall ./cmd/mall
.PHONY: mall

mall-arm: $(BUILD_DEPS)
	rm -f mall-arm
	GOOS=linux GOARCH=arm $(GOCC) build $(GOFLAGS) -o mall-arm ./cmd/mall
.PHONY: mall-arm

api-gen:
	$(GOCC) run ./gen/api
	goimports -w api
.PHONY: api-gen

cfgdoc-gen:
	$(GOCC) run ./node/config/cfgdocgen > ./node/config/doc_gen.go

build: mall 
.PHONY: build

install: install-mall 

install-mall:
	install -C ./mall /usr/local/bin/mall

mall-image:
	docker build -t mall:latest -f ./cmd/mall/Dockerfile .
.PHONY: mall-image

