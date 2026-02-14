# Req Client BDD

## Scope
Req-based runtime client for polling, claim, prompt resolution, state updates, and chat append.

## Scenarios
- Given runtime API key configuration, when polling for work, then eligible stories are returned.
- Given a story id, when claiming, then claim succeeds or returns conflict.
- Given project prompt preview request, when inputs are unchanged, then prompt response is deterministic.
- Given runtime chat/state events, when posting to API, then server persists updates.

## Invariants
- Use `Req` for all client HTTP calls.
- Do not include server-only concerns (router/db/auth plug logic) in client modules.
