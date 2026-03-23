BINARY     := valkyrie
CONFIG     ?= valkyrie.toml

.PHONY: build help

## build: compile the binary
build:
	go build -o $(BINARY) .

## help: show available targets
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
