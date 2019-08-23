COMMIT=$(shell git rev-parse HEAD | head -c 8)$(shell if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi)
REGISTRY?=docker.io
NAMESPACE?=stellarproject
IMAGE_NAMESPACE?=$(NAMESPACE)
APP=atlas
CLI=actl
REPO?=$(NAMESPACE)/$(APP)
TAG?=dev
BUILD?=-dev
BUILD_ARGS?=
PACKAGES=$(shell go list ./... | grep -v -e /vendor/)
CYCLO_PACKAGES=$(shell go list ./... | grep -v /vendor/ | sed "s/github.com\/$(NAMESPACE)\/$(APP)\///g" | tail -n +2)
VAB_ARGS?=
CWD=$(PWD)
BINARY_SUFFIX?=

ifneq "$(strip $(shell command -v go 2>/dev/null))" ""
	GOOS ?= $(shell go env GOOS)
	GOARCH ?= $(shell go env GOARCH)
else
	ifeq ($(GOOS),)
		# approximate GOOS for the platform if we don't have Go and GOOS isn't
		ifeq ($(OS),Windows_NT)
			GOOS = windows
		else
			UNAME_S := $(shell uname -s)
			ifeq ($(UNAME_S),Linux)
				GOOS = linux
			endif
			ifeq ($(UNAME_S),Darwin)
				GOOS = darwin
			endif
			ifeq ($(UNAME_S),FreeBSD)
				GOOS = freebsd
			endif
		endif
	else
		GOOS ?= $$GOOS
		GOARCH ?= $$GOARCH
	endif
endif

ifeq ($(GOOS),windows)
	BINARY_SUFFIX=".exe"
endif

all: binaries

generate:
	@>&2 echo " -> building protobufs for grpc"
	@echo ${PACKAGES} | xargs protobuild -quiet

bindir:
	@mkdir -p bin

binaries: cli daemon

cli: bindir
	@>&2 echo " -> building cli ${COMMIT}${BUILD}"
	@cd cmd/$(CLI) && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -installsuffix cgo -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT) -X github.com/$(REPO)/version.Build=$(BUILD)" -o ../../bin/$(CLI)$(BINARY_SUFFIX) .

daemon: bindir
	@>&2 echo " -> building daemon ${COMMIT}${BUILD}"
	@cd cmd/$(APP) && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -installsuffix cgo -ldflags "-w -X github.com/$(REPO)/version.GitCommit=$(COMMIT) -X github.com/$(REPO)/version.Build=$(BUILD)" -o ../../bin/$(APP)$(BINARY_SUFFIX) .

vet:
	@echo " -> $@"
	@test -z "$$(go vet ${PACKAGES} 2>&1 | tee /dev/stderr)"

lint:
	@echo " -> $@"
	@golint -set_exit_status ${PACKAGES}

check: vet lint

test:
	@go test -short -v -cover $(TEST_ARGS) ${PACKAGES}

install:
	@install -D -m 755 bin/* /usr/local/bin/

clean:
	@rm -rf bin/
	@rm -rf *.tar.gz
	@rm -rf build/

.PHONY: generate clean check test install binaries cli daemon lint
