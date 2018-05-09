# Utility script to build, test and push image to registry
CWD := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

VENDORBIN = $(CURDIR)/vendor/bin
PATH     := $(VENDORBIN):$(PATH)

KEYGEN_VERSION ?= $(shell git describe --tags 2>/dev/null ||  git rev-parse HEAD)
KEYGEN_REPO = quay.io/gravitational/keygen
BUILDBOX_TAG ?= golang:1.9.0-stretch
BUILDDIR ?=
export

.PHONY: all
all:
	mkdir -p build
	go build -o build/keygen github.com/gravitational/keygen/tool/keygen


# Build docker image
.PHONY: image
image: build-binary
	docker build \
           -t "$(KEYGEN_REPO):$(KEYGEN_VERSION)" .

# builds program inside Docker container
.PHONY: build-binary
build-binary:
	rm -rf $(CWD)/build
	mkdir -p $(CWD)/build
	docker run -v $(CWD)/build:/build -v $(CWD):/go/src/github.com/gravitational/keygen $(BUILDBOX_TAG) go build -o /build/keygen github.com/gravitational/keygen/tool/keygen


# Build chart
.PHONY: chart
chart:
	helm package keygen --version=$(KEYGEN_VERSION) -d $(BUILDDIR)

# install-chart
.PHONY: 
install-chart:
	helm install keygen --version=$(KEYGEN_VERSION) -n keygen

# upgrade chart
.PHONY: 
upgrade-chart:
	helm upgrade keygen keygen --version=$(KEYGEN_VERSION)

# Publish docker image. User runs this has to have Quay write permission
.PHONY: push
push:
	docker push $(KEYGEN_REPO):$(KEYGEN_VERSION)




