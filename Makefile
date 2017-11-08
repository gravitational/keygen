# Utility script to build, test and push image to registry
CWD := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

VENDORBIN = $(CURDIR)/vendor/bin
PATH     := $(VENDORBIN):$(PATH)

KEYGEN_VERSION ?= $(shell git describe --tags 2>/dev/null ||  git rev-parse HEAD)
KEYGEN_REPO = quay.io/gravitational/keygen
BUILDBOX_TAG ?= golang:1.9.0-stretch

# Build docker image
.PHONY: build
build: build-binary
	docker build \
           -t "$(KEYGEN_REPO):$(KEYGEN_VERSION)" .


# builds program inside Docker container
.PHONY: build-binary
build-binary:
	rm -rf $(CWD)/build
	mkdir -p $(CWD)/build
	docker run -v $(CWD)/build:/build -v $(CWD):/go/src/github.com/gravitational/keygen $(BUILDBOX_TAG) go build -o /build/keygen github.com/gravitational/keygen/tool/keygen

# Publish docker image. User runs this has to have Quay write permission
.PHONY: publishr
publish:
	docker push $(KEYGEN_REPO):$(KEYGEN_VERSION)


.PHONY: all
all:
	mkdir -p build
	go build -o build/keygen github.com/gravitational/keygen/tool/keygen


