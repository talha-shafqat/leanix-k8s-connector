GOARCH ?= amd64
GOOS ?= linux
GOVERSION ?= 1.12
PROJECT ?= leanix-k8s-connector
DOCKER_NAMESPACE ?= leanix

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

IMAGE := $(DOCKER_NAMESPACE)/$(PROJECT):$(VERSION)
BUILD_CMD=go build -o bin/$(PROJECT) -ldflags '-X $(go list -m)/pkg/version.VERSION=${VERSION} -extldflags "-static"' ./cmd/$(PROJECT)/main.go

TEST_CMD=go test ./pkg/...

DOCKER_CMD=docker run \
		--rm \
		--name $(PROJECT)-make \
		-e GOARCH=$(GOARCH) \
		-e GOOS=$(GOOS) \
		-e CGO_ENABLED=0 \
		-v $(PWD):/tmp/$(PROJECT) \
		-w /tmp/$(PROJECT) \
		golang:$(GOVERSION) \

DOCKER_BUILD_CMD=$(DOCKER_CMD) \
		$(BUILD_CMD)

DOCKER_TEST_CMD=$(DOCKER_CMD) \
		$(TEST_CMD)

ifdef GOPATH
DOCKER_BUILD_CMD=docker run \
		--rm \
		--name $(PROJECT)-make \
		-e GOARCH=$(GOARCH) \
		-e GOOS=$(GOOS) \
		-e CGO_ENABLED=0 \
		-v $(GOPATH)/pkg:/go/pkg \
		-v $(PWD):/tmp/$(PROJECT) \
		-w /tmp/$(PROJECT) \
		golang:$(GOVERSION) \
		$(BUILD_CMD)

DOCKER_TEST_CMD=docker run \
		--rm \
		--name $(PROJECT)-make \
		-e GOARCH=$(GOARCH) \
		-e GOOS=$(GOOS) \
		-e CGO_ENABLED=0 \
		-v $(GOPATH)/pkg:/go/pkg \
		-v $(PWD):/tmp/$(PROJECT) \
		-w /tmp/$(PROJECT) \
		golang:$(GOVERSION) \
		$(TEST_CMD)
endif

.PHONY: all

all: clean test build

local: clean test-local build-local

clean:
	$(RM) bin/$(PROJECT)

build:
	$(DOCKER_BUILD_CMD)

build-local:
	$(BUILD_CMD)

image:
	docker build -t $(IMAGE) .

push:
	docker push $(IMAGE)

test:
	$(DOCKER_TEST_CMD)

test-local:
	$(TEST_CMD)