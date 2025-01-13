help:
	@echo ""
	@echo " ########################"
	@echo " # usage: make <target> #"
	@echo " ########################"
	@echo ""
	@echo " available targets:"
	@echo "   - build      -> Builds the binary in 'bin/dotcomfy'"
	@echo "   - references -> Scrapes citations in code and builds References.md"
	@echo ""

.PHONY: help

build:
	go build -o bin/dotcomfy main.go

.PHONY: build

references:
	./scripts/buildRefs.sh
	@echo "Built docs/REFERENCES.md"

.PHONY: references

IMAGE_NAME := dotcomfy
IMAGE_TAG := latest

CONTAINER_RUNTIME := $(shell command -v podman >/dev/null 2>&1 && echo podman || echo docker)

build-container:
	$(CONTAINER_RUNTIME) build -t $(IMAGE_NAME):$(IMAGE_TAG) -f Containerfile

.PHONY: build-container

container: build-container
	$(CONTAINER_RUNTIME) run --rm -it $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: container
