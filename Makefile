GOOS ?= $(shell go env GOOS)
GIT_VERSION := $(shell git describe --always --tags)
GIT_BRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)
GIT_HASH := $(GIT_BRANCH)/$(shell git log -1 --pretty=format:"%H")

REGISTRY?=yashvardhankukreja
REPO=$(REGISTRY)/kyverno
IMAGE_TAG?=$(GIT_VERSION)

IMAGE_NAME := kube-bench-exporter

docker-build-dev: fmt tidy vet
	@docker build -t $(REPO)/$(IMAGE_NAME):$(IMAGE_TAG) .

tidy:
	go mod tidy

vet:
	go vet ./...

# Run go fmt against code
fmt:
	gofmt -s -w .

fmt-check:
	gofmt -s -l .

