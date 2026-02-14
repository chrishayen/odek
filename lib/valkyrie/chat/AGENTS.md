# Chat Domain BDD

## Scope
Project-scoped chat persistence and realtime project event delivery.

## Scenarios
- Given project chat thread creation with optional story link, when creating threads, then threads persist as project-scoped and story link is optional.
- Given frontend and runtime messages through API proxy, when messages are created, then all messages are persisted and list endpoints return full history.
- Given SSE subscription for a project, when a new message is created in that project, then a chat message event is emitted.
- Given SSE subscription for project A, when messages are created in project B, then no events are emitted to project A stream.

## Invariants
- Chat threads are project-scoped.
- Realtime stream transport is SSE in v1.
