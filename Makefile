BINARY     := valkyrie
CONFIG     ?= valkyrie.toml
DEV_DIR    := .dev
PORT       := 8319

.PHONY: build install run tui decompose hydrate generate serve test test-e2e dev clean help

## build: compile the binary
build:
	go build -o $(BINARY) .

## install: compile and install to GOBIN
install:
	go install .

## tui: install and launch the TUI
tui: install
	@valkyrie tui

## run: build and run with ARGS
run: build
	@./$(BINARY) $(ARGS)

## decompose: decompose requirements (pass ARGS="your requirement text")
decompose: build
	@./$(BINARY) runes decompose $(ARGS)

## hydrate: hydrate all un-hydrated runes
hydrate: build
	@./$(BINARY) runes hydrate-all $(ARGS)

## generate: decompose then hydrate (full pipeline)
generate: decompose
	@./$(BINARY) runes hydrate-all

## serve: start the HTTP API server
serve: build
	@./$(BINARY) serve --port $(PORT)

## test: run all tests
test:
	go test ./e2e/ -v

## test-e2e: run end-to-end tests
test-e2e:
	go test ./e2e/ -v

## dev: install, reset .dev/, init, and run claude
dev: install
	rm -rf $(DEV_DIR)
	mkdir -p $(DEV_DIR)
	cd $(DEV_DIR) && valkyrie init go
	cd $(DEV_DIR) && claude "I want to create a cli that starts an http server and serves the directory where it was started"

## clean: remove binary and generated code
clean:
	rm -f $(BINARY)
	rm -rf src/

## help: show available targets
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
