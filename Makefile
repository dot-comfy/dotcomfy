help:
	@echo ""
	@echo " ########################"
	@echo " # usage: make <target> #"
	@echo " ########################"
	@echo ""
	@echo " available targets:"
	@echo "   - build -> Builds the binary in 'bin/dotcomfy'"
	@echo ""

.PHONY: help

build:
	go build -o bin/dotcomfy cmd/dotcomfy/main.go

.PHONY: build
