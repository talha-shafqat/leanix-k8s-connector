PROJECT ?= leanix-k8s-connector
DOCKER_NAMESPACE ?= leanix

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)
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
	CGO_ENABLED=0 go build -o bin/$(PROJECT) -ldflags '-X $(go list -m)/pkg/version.VERSION=${VERSION} -extldflags "-static"' ./cmd/$(PROJECT)/main.go

version:
	@echo $(VERSION)

image:
	docker build -t $(IMAGE) $(LATEST) .

push:
	docker push $(IMAGE) $(LATEST)

test:
	go test ./pkg/...
