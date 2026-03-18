# Valkyrie — Agentic Code Server (Draft Spec)

> A code infrastructure system designed for agentic development, not human development.

## Problem

Standard software development patterns break down with AI agents:
- Agents can't hold large codebases in context
- Agents code faster than humans can review
- Agents are sloppy — mistakes compound without guardrails
- Monorepos are too cumbersome; per-feature repos are too fragmented
- File system is a terrible interface for agents

## Core Idea

Replace direct filesystem access with a **Code Server** — a semantic interface between agents and code.

Instead of reading files, agents query the server:
- "What handles user authentication?"
- "What would break if I change this contract?"
- "Here's my change — does it satisfy requirement X?"

The server is the source of truth about what the system is.

## The Rune

The atomic unit of functionality. Not a file, not a repo, not a service — a **named, registered piece of functionality** with a contract.

A rune has a specific shape and means a specific thing. Like its namesake, it's a symbol with precise, defined meaning — carve it, and it means exactly one thing. Combine runes and you cast something larger.

An agent gets assigned a rune (or a task scoped to one), queries the server for context, does the work, hands it back. The server validates fit before accepting it.

### Rune Record (rough shape)

```
name: string                   # unique identifier
description: text              # what it does (human + agent readable)
version: semver                # current version
stage: draft | reviewed | stable   # promotion state
inputs: schema ref             # what it accepts
outputs: schema ref            # what it returns
events_published: [event ref]  # broadcasts
events_subscribed: [event ref] # listens to
dependencies: [rune ref]       # other runes it depends on
requirements: [req ref]        # requirements this rune satisfies
config: [config key ref]       # runtime config it needs
```

## Wiring: Broadcast / Respond

Runes communicate via events, not direct calls.

- Loose coupling — runes don't know about each other, only message contracts
- Server is the broker — sees all traffic, can audit, validate, monitor health
- Adding a new rune that reacts to existing events requires zero changes elsewhere
- Server maintains a full subscription map — agents can query "what does this event trigger?"

Text is the source of truth for wiring. Graphs are a rendered view for humans, generated on demand.

## Registries (what the server tracks)

Inspired by: Schema Registry (compatibility enforcement), Consul (discovery + health), npm (manifest + deps), MLflow (promotion stages).

| Registry | What it tracks |
|---|---|
| Rune Registry | Named units of functionality, their state and stage |
| Schema Registry | Data shapes — DB tables, types, event payloads |
| API Registry | External-facing contracts — endpoints, inputs, outputs, errors |
| Event Registry | Broadcastable messages — name, payload shape, publisher |
| Subscription Registry | Who listens to what events |
| Requirement Registry | User requirements mapped to runes that satisfy them |
| Dependency Registry | Rune → rune dependency graph |
| Config Registry | Per-rune runtime config (no more scattered .env files) |

All entries share a common shape: `name`, `description`, `version/state`, `relationships`.

## Storage

Plain text files — TOML per entry, one file per rune/event/schema, flat directory structure. Human-readable, git-diffable, no database required.

```
registry/
  runes/
    user-auth.toml
    payment-processor.toml
  events/
    user-signed-up.toml
  schemas/
    user.toml
```

The server loads the registry into memory on start. Git handles history.

## Protocol

Valkyrie exposes an **MCP (Model Context Protocol)** server. Any MCP-compatible agent can connect and use it without pre-configuration — the server declares its tools and the agent discovers them dynamically.

Built in **Go** — compiles to a single binary, easy to install on any platform.

## Human Review Model

Instead of reviewing code diffs, humans review **contract changes**:
- Schema changed
- API contract changed
- Rune X now does Y instead of Z
- Rune promoted from `reviewed` → `stable`

**Promotion stages** (from ML model registry pattern):
`draft` → `reviewed` → `stable`

This is the human checkpoint — not a PR, not a diff. A promotion approval.

## Agent Interaction Model

1. Agent receives a task scoped to a rune (or a new rune to create)
2. Agent queries server: "give me context for rune X" → server returns contract, deps, relevant schemas, config keys
3. Agent does the work
4. Agent submits change to server
5. Server validates: contracts satisfied? dependencies intact? requirements met? breaking changes flagged?
6. If valid: rune advances to `draft` or `reviewed` depending on confidence
7. Human approves promotion to `stable`

## Config / Secrets

Per-rune config owned by the server. Agent asks: "what does this rune need to run?" — no digging through env files. Server returns config keys; actual secrets resolved at runtime.

## Open Questions

- [ ] Write path: how does the server accept code changes? Diff-based? AST? Full replacement?
- [ ] What does the agent-facing MCP tool surface look like exactly?
- [ ] Wiring config format — declarative? Who can propose changes?
- [ ] How do runes get physically stored/deployed?
- [ ] What's the minimum viable version we could build and use ourselves first?

## Status

Early spec — brainstormed with Chris 2026-03-18.
