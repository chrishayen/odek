# Odek — The Rune Server

> Describe functionality in English. Odek builds it, registers it, and wires it together — in any language, in parallel.

---

## The Problem

Standard software development patterns break down with AI agents:
- Agents can't hold large codebases in context
- Agents code faster than humans can review
- Agents are sloppy — mistakes compound without guardrails
- Monorepos are too cumbersome; per-feature repos are too fragmented
- The filesystem is a terrible interface for agents

## The Idea

Odek is a **rune server** — an orchestration layer between you and your codebase.

You describe what you need. Odek spins up sandboxed agents to build it, registers the result, and tracks how everything connects. Agents work in parallel. Humans review contracts, not diffs.

## The Rune

The atomic unit of functionality. Not a file, not a repo, not a service — a **named, registered piece of functionality** with a defined contract.

A rune has a specific shape and means a specific thing. Like its namesake, it's a symbol with precise meaning — carve it, and it means exactly one thing. Combine runes and you cast something larger.

### Rune Record

Runes are stored as markdown with YAML frontmatter:

```markdown
---
version: 0.1.0
hydrated: false
coverage: -1
---

# auth/validate-email

Validates that a string is a well-formed email address.

## Signature

(email: string) -> bool

## Behavior

- Input: a string
- Output: boolean
- Returns true if string contains exactly one @ with non-empty local and domain parts
- Returns false otherwise

## Positive tests

- Given 'user@example.com', returns true
- Given 'a@b.co', returns true

## Negative tests

- Given an empty string, returns false
- Given 'no-at-sign', returns false
```

### Signature Types

Signatures use precise types with numeric precisions:

| Category | Types |
|---|---|
| Signed integers | `i8`, `i16`, `i32`, `i64` |
| Unsigned integers | `u8`, `u16`, `u32`, `u64` |
| Floating point | `f32`, `f64` |
| Other primitives | `string`, `bool`, `bytes` |
| Collections | `list[T]`, `map[K, V]` |
| Nullable | `optional[T]` |
| Fallible | `result[T, E]` |

Use `result[T, E]` for any function that can fail. Types nest: `result[list[i32], string]`.

## How It Works

1. You describe requirements in English
2. Odek's analyzer decomposes them into runes — one function each
3. You review and approve the proposed runes (name, description, signature, behavior, tests)
4. Odek registers each rune in the registry
5. The hydrator spins up a sandboxed agent to generate code + tests for a rune
6. Odek runs the tests, tracks coverage, and marks the rune as hydrated

### Analyze

```bash
odek runes analyze --requirements "User login with email and password"
```

The analyzer breaks requirements into the smallest functional parts, checks existing runes to avoid duplication, and proposes new runes with signatures, behavior specs, and test cases.

### Hydrate

```bash
odek runes hydrate auth/validate-email
```

The hydrator sends the rune spec to a sandboxed coding agent, extracts the generated files, runs tests, and records coverage.

## The Feature

Runes are atomic — one function each. A **feature** groups runes into something larger. It describes the namespace, its components, and how runes wire together.

A feature lives in the same directory as its runes: `runes/auth/feature.md` sits alongside `runes/auth/validate-email.md`.

### Feature Record

```markdown
---
version: 0.1.0
status: draft
---

# auth

User authentication and authorization.

## Components

### login

Accepts a username and password, validates both, returns a session token.

Composes: auth/validate-email, auth/hash-password, auth/create-session-token

### password-reset

Accepts an email, verifies the account exists, sends a reset link.

Composes: auth/validate-email, auth/lookup-account-by-email, notifications/send-email

## Connections

- login receives (username: string, password: string) -> result[string, string]
- password-reset receives (email: string) -> result[bool, string]
- login output feeds into session/store-session
```

Components describe how runes compose. The `Composes:` line lists runes by name. Connections describe data flow between components and other features.

## CLI

```bash
odek init                          # Initialize a new project
odek runes list                    # List all runes
odek runes create --name <slug> --description <desc> --signature <sig>
odek runes get <name>              # Get a rune by name
odek runes update <name> --description <desc> --signature <sig> --version <ver>
odek runes delete <name>           # Delete a rune
odek runes analyze --requirements <text> [--yes]
odek runes hydrate <name>          # Generate code via sandbox agent
odek features list                 # List all features
odek features create --name <ns> --description <desc>
odek features get <name>           # Get a feature by name
odek features update <name> --description <desc> --version <ver> --status <status>
odek features delete <name>        # Delete a feature
```

## MCP Server

Odek exposes an **MCP (Model Context Protocol)** server. Any MCP-compatible agent connects and discovers tools dynamically.

```bash
odek mcp    # Start MCP server (stdio transport)
```

Tools: `runes_list`, `runes_create`, `runes_get`, `runes_update`, `runes_delete`, `runes_analyze`, `runes_create_batch`, `runes_hydrate`, `features_list`, `features_create`, `features_get`, `features_update`, `features_delete`.

Running `odek init` generates `.mcp.json` so Claude Code auto-discovers the server.

## Storage

Plain markdown files. Human-readable, git-diffable, no database.

```
runes/
  auth/
    feature.md              # feature description
    validate-email.md       # rune spec
    validate-email/         # generated code (after hydration)
      main.go
      main_test.go
  payment/
    feature.md
    calculate-total.md
    calculate-total/
```

Namespaces are organizational (e.g. `auth/`, `payment/`), not judgments about purity.

## Configuration

```toml
# odek.toml
project = "my-app"

[agent]
type = "claude-sub"              # "claude-sub" (default) or "mock" (testing)
model = "claude-sonnet-4-5"      # optional model override
token = "sk-ant-..."             # optional: literal API token
token_env = "MY_TOKEN"           # optional: env var name for token
image = "ghcr.io/..."            # optional: custom Docker image
```

## Protocol

Built in **Go** — single binary, runs on Linux/Mac/Windows.

The sandbox agent runs inside Docker. Token resolution order:
1. `token` field in config
2. `token_env` field in config
3. `CLAUDE_CODE_OAUTH_TOKEN` env var
4. `~/.claude/.credentials.json`

## Roadmap

**Done**
- [x] Rune registry — markdown files, git-backed
- [x] CLI — create, get, list, update, delete
- [x] MCP server — exposes registry to agents
- [x] Analyze — English requirements → rune decomposition via sandbox agent
- [x] Hydrate — rune spec → code + tests via sandbox agent
- [x] Coverage tracking
- [x] Function signatures with type precisions

**Next**
- [x] Wiring — features describe how runes compose into components
- [ ] Validation — contract checks on submission
- [ ] Promotion workflow — draft → reviewed → stable
- [ ] Parallel hydration — multiple runes built simultaneously

**Future**
- [ ] Event system — broadcast/respond between runes
- [ ] Schema registry — data shapes, DB tables, event payloads
- [ ] Hosted sandbox execution
- [ ] Public rune registry — publish and pull shared runes
