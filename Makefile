APP_NAME?=geogame
IMAGE=$(APP_NAME)
VERSION?=snap-shot
GOCMD=CGO_ENABLED=0 GOOS=linux go

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
.PHONY: build
build:
	$(GOCMD) build -mod vendor -ldflags "-X main.serviceVersion=$(VERSION)" -o go-app $(ROOT_DIR)

.PHONY: image
image:
	@echo "Building $(IMAGE):$(VERSION)"
	@docker build -f Dockerfile --build-arg=VERSION=$(VERSION) -t $(IMAGE):$(VERSION) .