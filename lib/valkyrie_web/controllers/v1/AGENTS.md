# V1 Controller Contract BDD

## Scope
Public and protected HTTP contracts for v1 JSON API.

## Scenarios
- Given v1 resource catalog requests, when listing resources, then unsupported resources like `agents`, `memory_entries`, and `context_entries` are absent.
- Given bootstrap user creation, when a user is created through `/v1/users`, then the response includes the resolved `organization_id` and role used for tenancy-scoped access.
- Given protected endpoints, when authentication fails, then JSON error envelope is returned with auth failure code.
- Given authorization failures, when role lacks permission, then API returns forbidden.
- Given soft-deleted entities, when reading default endpoints, then API returns not found.
- Given list endpoints, when cursor is supplied, then pagination continues without duplicate items.

## Invariants
- API errors use `%{error: %{code: ..., message: ...}}`.
- All protected controller actions enforce organization scoping.
