REPO_NAME := skeleton-go-api
BINARY_NAME := $(REPO_NAME)
# Fetch the latest git tag.
GIT_TAG := $(shell git describe --tags `git rev-list --tags --max-count=1`)
# Use the latest git tag as the image tag. If no tag is found, use "latest".
IMAGE_TAG := $(if $(GIT_TAG),$(GIT_TAG),latest)

# Default target (since it's the first without '.' prefix)
build-all: coverage build
.PHONY: build-all

build:
	go build ./cmd/$(BINARY_NAME)
.PHONY: build

test:
	go test -race -v ./...
.PHONY: test

coverage:
	./script/coverage.sh
.PHONY: cover

lint:
	golangci-lint run -v
.PHONY: lint

lint-docker:
	docker run --rm -e CGO_ENABLED=0 -v $(pwd):/app -w /app golangci/golangci-lint:v1.57.1 golangci-lint run -v
.PHONY: lint-docker

docker-build:
	docker build --tag "$(BINARY_NAME):$(IMAGE_TAG)" .
.PHONY: docker-build
