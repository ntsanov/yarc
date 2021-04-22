SHELL := /bin/bash
ifeq ($(VERSION),)
	VERSION := $(shell git rev-list --count HEAD)
endif
ifeq ($(COMMIT),)
	COMMIT := $(shell git describe --always --long --dirty)
endif
BUILDTIME=$(shell date -u "+%Y-%m-%dT%H:%M:%S")
env := GO111MODULE=on
ldflags := -X main.version=v${VERSION} -X main.commit=${COMMIT} -X main.buildtime=${BUILDTIME}
ldflags_static := $(ldflags) -w -extldflags '-static'
cli = ./build/yarc

all:
	go build -o $(cli) -ldflags="$(ldflags)"

static:
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/mcl
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/bls minimised_static BLS_SWAP_G=1
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build -o $(cli) -ldflags="$(ldflags_static)" main.go

clean:
	go clean
	rm -rf ./build
