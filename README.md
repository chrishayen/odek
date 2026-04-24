# ODEK

ODEK is a terminal app for decomposing software requirements into hierarchical
function trees called runes. It talks to an OpenAI-compatible chat completions
API, reads a local corpus of worked examples, and renders the result in a
Bubble Tea TUI.

## Requirements

- Go 1.26.1 or newer
- An OpenAI-compatible API server with `/v1/chat/completions`
- A terminal with true color support

By default ODEK connects to `http://localhost:8080/v1`. Set `API_BASE_URL` to
point somewhere else. You can provide credentials with `OPENAI_API_KEY` or
`API_KEY`.

`API_BASE_URL` can include or omit `/v1`; ODEK normalizes either form.

## Run

Start the TUI:

```sh
go run .
```

Direct chat:

```sh
go run . -p "Explain what ODEK does"
```

Direct chat as raw JSON:

```sh
go run . -p "Explain what ODEK does" -json
```

Direct decomposition as JSON:

```sh
go run . -d "a JWT signer and verifier"
```

## Configuration

- `API_BASE_URL`: API origin or versioned base URL. Defaults to
  `http://localhost:8080`.
- `OPENAI_API_KEY`: bearer token for hosted APIs.
- `API_KEY`: fallback bearer token when `OPENAI_API_KEY` is unset.

## Validate

```sh
go build ./...
go test ./...
```

The GitHub Actions workflow runs both commands for pushes and pull requests to
`master`.

## Repo Map

- `main.go`: CLI flags, direct modes, and TUI startup.
- `openai/`: small OpenAI-compatible client with streaming and tool-call loops.
- `internal/decomposer/`: embedded prompt, example retrieval, decomposition,
  merge, and recursive expansion engine.
- `internal/tui/`: landing screen, chat, split pane, and rune browser.
- `internal/examples/`: loader and lookup index for the example corpus.
- `examples/`: worked decomposition examples used by the model.
- `PRD.md`: product and architecture reference.
