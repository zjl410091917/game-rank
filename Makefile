GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

PROTO_FILES=$(shell find -L protos -name *.proto)
IPROTO_FILES=$(shell find iprotos -name *.proto)

GIT_BRANCH=$(shell git symbolic-ref --short HEAD)
GIT_CMID=$(shell git rev-parse --short HEAD)
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy mod file
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	golangci-lint run -v
	go mod verify


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the cmd/api application
.PHONY: build
build: audit
	go build -ldflags='-s -w' -o=./bin/server ./cmd

.PHONY: sim
sim:
	go build -ldflags='-s -w' -o=./bin/sim_us ./sim/sim_update_score

