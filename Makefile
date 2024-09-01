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
	go build -o bin/dotcomfy cmd/dotcomfy/main.go

.PHONY: build

references:
	./scripts/buildRefs.sh
	@echo "Built References.md"

.PHONY: references
