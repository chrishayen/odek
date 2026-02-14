# API Keys Domain BDD

## Scope
API key authentication lifecycle and organization-scoped runtime access.

## Scenarios
- Given a user membership in an organization, when creating an API key, then the raw key is returned once and metadata is persisted.
- Given active API keys, when listing keys, then metadata is returned and raw secrets are never returned.
- Given an active key, when revoking the key, then future authentication with that key fails.
- Given a key scoped to organization A, when calling organization B resources, then access is denied.

## Invariants
- Never store raw API key secrets.
- Store key hash + non-sensitive prefix for lookup/debug metadata.
- API key authentication must load organization, user, role, and key id into request scope.
