BINARY     := valkyrie
CONFIG     ?= valkyrie.toml
DEV_DIR    := .dev

.PHONY: build install dev help

## build: compile the binary
build:
	go build -o $(BINARY) .

## install: compile and install to GOBIN
install:
	go install .

## dev: install, reset .dev/, init, and run claude
dev: install
	rm -rf $(DEV_DIR)
	mkdir -p $(DEV_DIR)
	cd $(DEV_DIR) && valkyrie init
	cd $(DEV_DIR) && claude "I want to create a cli that starts an http server and serves the directory where it was started"

## help: show available targets
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
