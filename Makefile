BINARY     := valkyrie
CONFIG     ?= valkyrie.toml
BASE_URL   ?= http://localhost:8080

.PHONY: build dev rune getrune help

## build: compile the binary
build:
	go build -o $(BINARY) .

## dev: start the server (requires $(CONFIG))
dev: build
	VALKYRIE_CONFIG=$(CONFIG) ./$(BINARY) serve

## rune: create a rune — usage: make rune DESCRIPTION="Returns Hello World"
rune:
ifndef DESCRIPTION
	$(error DESCRIPTION is required. Usage: make rune DESCRIPTION="your description here")
endif
	@NAME=$$(echo "$(DESCRIPTION)" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '-' | sed 's/^-//;s/-*$$//'); \
	echo "Creating rune: $$NAME"; \
	curl -s -X POST $(BASE_URL)/runes \
		-H "Content-Type: application/json" \
		-d "{\"name\":\"$$NAME\",\"description\":\"$(DESCRIPTION)\"}" | python3 -m json.tool

## getrune: get a rune by name — usage: make getrune NAME=hello-world
getrune:
ifndef NAME
	$(error NAME is required. Usage: make getrune NAME=hello-world)
endif
	@curl -s $(BASE_URL)/runes/$(NAME) | python3 -m json.tool

## help: show available targets
help:
	@grep -E '^## ' Makefile | sed 's/^## //'
