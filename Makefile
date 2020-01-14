PROJECT ?= leanix-k8s-connector
DOCKER_NAMESPACE ?= leanix

# This version-strategy uses git tags to set the version string
VERSION := 2.0.0-beta1-$(shell git describe --tags --always --dirty)
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

IMAGE := $(DOCKER_NAMESPACE)/$(PROJECT):$(VERSION)
LATEST := $(DOCKER_NAMESPACE)/$(PROJECT):latest

.PHONY: all

all: clean test build

clean:
	$(RM) bin/$(PROJECT)

build:
	CGO_ENABLED=0 go build -o bin/$(PROJECT) -ldflags '-X $(shell go list -m)/pkg/version.VERSION=${VERSION} -extldflags "-static"' ./cmd/$(PROJECT)/main.go

version:
	@echo $(VERSION)

image:
	docker build -t $(IMAGE) -t $(LATEST) .

push:
	docker push $(IMAGE)
	docker push $(LATEST)

test:
	go test ./pkg/...
