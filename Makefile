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

IMAGE=fedora:latest
CONTAINER_NAME=temp-fedora
SKEL_DIR=~

# The following two targets get the contents of a typical Linux user's home
# directory to use as a base for our tests

pull-image:
	@echo "Pulling $(IMAGE)"
	@$(CONTAINER_RUNTIME) pull $(IMAGE)

.PHONY: pull-image

extract-home:
	@echo "Creating temporary container..."
	$(CONTAINER_RUNTIME) run -d --name $(CONTAINER_NAME) $(IMAGE) sleep infinity
	@echo "Copying /etc/skel to $(SKEL_DIR)"
	$(CONTAINER_RUNTIME) cp $(CONTAINER_NAME):/etc/skel $(SKEL_DIR)
	$(CONTAINER_RUNTIME) stop $(CONTAINER_NAME)
	$(CONTAINER_RUNTIME) rm $(CONTAINER_NAME)

.PHONY: extract-home

build-container:
	@$(MAKE) pull-image
	@$(MAKE) extract-home
	$(CONTAINER_RUNTIME) build -t $(IMAGE_NAME):$(IMAGE_TAG) -f Containerfile

.PHONY: build-container

container: build-container
	$(CONTAINER_RUNTIME) run --rm -it $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: container

TEST_DIR := tests/scripts

check-test:
	@if [ ! -f $(TEST_DIR)/$(TEST_SCRIPT) ]; then \
		echo "Error: Test script $(TEST_DIR)/$(TEST_SCRIPT) not found!"; \
		exit 1; \
	fi

.PHONY: check-test

.PHONY: test-%

# Running `test-install` will run `tests/scripts/install.sh` from inside the
# container.
test-%:
	$(MAKE) TEST_SCRIPT=$*.sh check-test build-container
	@echo "Running test $*.sh in container ..."
	$(CONTAINER_RUNTIME) run --rm $(IMAGE_NAME):$(IMAGE_TAG) bash $(TEST_DIR)/$*.sh

# This runs all the test scripts
test:
	@for test in $(wildcard $(TEST_DIR)/*.sh); do \
		@echo $$test; \
		test_name=$$(basename $$test .sh); \
		$(MAKE) test-$$test_name; \
	done

.PHONY: test
