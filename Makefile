GOARCH=amd64
GOOS=linux
GOVERSION=1.12
PROJECT=leanix-k8s-connector

BUILD_CMD=docker run \
		--rm \
		--name $(PROJECT)-build \
		-e GOARCH=$(GOARCH) \
		-e GOOS=$(GOOS) \
		-v $(PWD):/tmp/$(PROJECT) \
		-w /tmp/$(PROJECT) \
		golang:$(GOVERSION) \
		go build

ifdef GOPATH
BUILD_CMD=docker run \
		--rm \
		--name $(PROJECT)-build \
		-e GOARCH=$(GOARCH) \
		-e GOOS=$(GOOS) \
		-v $(GOPATH)/pkg:/go/pkg \
		-v $(PWD):/tmp/$(PROJECT) \
		-w /tmp/$(PROJECT) \
		golang:$(GOVERSION) \
		go build
endif

.PHONY: all

all: clean build

clean:
	$(RM) $(PROJECT)

build:
	$(BUILD_CMD)
