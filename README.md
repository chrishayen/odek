# Valkyrie — The Rune Server

> Describe functionality in English. Valkyrie builds it, registers it, and wires it together — in any language, in parallel.

---

## The Problem

Standard software development patterns break down with AI agents:
- Agents can't hold large codebases in context
- Agents code faster than humans can review
- Agents are sloppy — mistakes compound without guardrails
- Monorepos are too cumbersome; per-feature repos are too fragmented
- The filesystem is a terrible interface for agents

## The Idea

Valkyrie is a **rune server** — an orchestration layer between you and your codebase.

You describe what you need. Valkyrie spins up sandboxed agents to build it, registers the result, and tracks how everything connects. Agents work in parallel. Humans review contracts, not diffs.

## The Rune

The atomic unit of functionality. Not a file, not a repo, not a service — a **named, registered piece of functionality** with a defined contract.

A rune has a specific shape and means a specific thing. Like its namesake, it's a symbol with precise meaning — carve it, and it means exactly one thing. Combine runes and you cast something larger.

### Rune Record

```toml
name = "user-auth"
description = "Handles user authentication via email/password and JWT tokens"
version = "0.1.0"
stage = "draft"            # draft | reviewed | stable
runtime = "go@1.22"        # optional — language hint for code generation
path = "runes/user-auth"   # where the code lives

inputs = ["credentials-schema"]
outputs = ["auth-token-schema"]
events_published = ["user-authenticated", "auth-failed"]
events_subscribed = []
dependencies = ["token-generator", "user-store"]
requirements = ["REQ-001", "REQ-002"]
config = ["JWT_SECRET", "TOKEN_TTL"]
```

## How It Works

1. You (or your agent) describe a rune in English
2. Valkyrie queues it and spins up a sandboxed coding agent
3. The agent generates code + tests to satisfy the contract
4. Valkyrie registers the rune, validates contracts, checks for breakage
5. Rune lands in `draft` — human promotes to `reviewed` → `stable`
6. Repeat for 20 runes simultaneously if needed

## Wiring: Broadcast / Respond

Runes communicate via events, not direct calls.

- Loose coupling — runes don't know about each other, only message contracts
- Valkyrie is the broker — sees all traffic, audits, validates, monitors
- Adding a rune that reacts to existing events requires zero changes elsewhere
- Full subscription map — query "what does this event trigger?" at any time

Wiring is stored as text. Graphs are rendered on demand for humans.

## Storage

Plain TOML files, one per rune/event/schema. Human-readable, git-diffable, no database.

```
registry/
  runes/
    user-auth.toml
    payment-processor.toml
  events/
    user-authenticated.toml
  schemas/
    user.toml
  requirements/
    REQ-001.toml
```

Valkyrie loads the registry into memory on start. Git handles history.

## Registries

| Registry | What it tracks |
|---|---|
| Rune Registry | Named units of functionality, stage, contracts |
| Schema Registry | Data shapes — DB tables, types, event payloads |
| API Registry | External-facing contracts — endpoints, inputs, outputs |
| Event Registry | Broadcastable messages — name, payload, publisher |
| Subscription Registry | Who listens to what |
| Requirement Registry | Requirements mapped to runes that satisfy them |
| Dependency Registry | Rune → rune dependency graph |
| Config Registry | Per-rune runtime config keys |

## Protocol

Valkyrie exposes an **MCP (Model Context Protocol)** server. Any MCP-compatible agent connects and discovers tools dynamically — no pre-configuration needed.

Built in **Go** — single binary, runs on Linux/Mac/Windows.

## Sandbox

Valkyrie runs coding agents in isolated sandboxes to generate rune implementations. Sandbox provider is pluggable:

```toml
# valkyrie.toml
sandbox = "docker"         # default: local Docker
# sandbox = "valkyrie-cloud"  # hosted: fast, managed, billed
```

## Human Review

Humans review **contract changes**, not code diffs:
- Schema changed
- API contract changed
- Rune X now does Y instead of Z
- Rune promoted from `reviewed` → `stable`

`draft` → `reviewed` → `stable`

## Roadmap

**MVP**
- [ ] Rune registry — TOML files, flat directory, git-backed
- [ ] CLI — `create`, `get`, `list`, `update` rune
- [ ] MCP server — exposes registry to agents
- [ ] Rune record schema — full field set

**Agent Orchestration**
- [ ] English description → sandboxed agent → code + tests → registered rune
- [ ] Parallel execution — multiple runes built simultaneously
- [ ] Pluggable sandbox provider (default: Docker)

**Validation & Review**
- [ ] Contract validation on submission — breaks deps? conflicts?
- [ ] Promotion workflow — `draft` → `reviewed` → `stable`
- [ ] Human approval gate

**Hosted / Cloud**
- [ ] Managed sandbox execution (`sandbox: valkyrie-cloud`)
- [ ] Public rune registry — publish and pull shared runes
- [ ] Billing

---

*Brainstormed with Chris — 2026-03-18*
