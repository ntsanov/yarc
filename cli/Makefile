SHELL := /bin/bash
BIN_NAME=yarc
LDFLAGS=-w -extldflags '-static'
ENV:=GO111MODULE=on 
BUILD_DIR=./build
CC=musl-gcc
env := GO111MODULE=on

all:
	go build -o $(BUILD_DIR)/$(BIN_NAME) main.go

static:
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/mcl
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/bls minimised_static BLS_SWAP_G=1
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build -o $(BUILD_DIR)/$(BIN_NAME) -ldflags="$(LDFLAGS)" main.go

clean:
	go clean
	rm -rf ${BUILD_DIR}
