BINARY     := valkyrie
CONFIG     ?= valkyrie.toml

.PHONY: build install help

## build: compile the binary
build:
	go build -o $(BINARY) .

## install: compile and install to GOBIN
install:
	go install .

## help: show available targets
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
