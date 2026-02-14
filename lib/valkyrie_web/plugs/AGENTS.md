# API Plug BDD

## Scope
API key authn/authz plug behavior for v1 protected routes.

## Scenarios
- Given missing/invalid bearer token, when protected routes are called, then authentication fails with unauthorized.
- Given revoked or unknown API key, when protected routes are called, then authentication fails.
- Given valid API key, when protected routes are called, then request assigns include principal user/org/role/key metadata.
- Given permission checks, when role is insufficient, then request halts with forbidden error envelope.

## Invariants
- API key auth is separate from browser session auth.
- Plugs should not leak raw key material into assigns or logs.
