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

TEST_DIR := tests/scripts

check-test:
	@if [ ! -f $(TEST_DIR)/$(TEST_SCRIPT) ]; then \
		echo "Error: Test script $(TEST_DIR)/$(TEST_SCRIPT) not found!"; \
		exit 1; \
	fi

.PHONY: check-test

BINDIR = /usr/local/bin
BINARY = ./bin/dotcomfy

install:
	@echo "Installing dotcomfy to $(BINDIR)"
	sudo install -m 755 $(BINARY) $(BINDIR)

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
