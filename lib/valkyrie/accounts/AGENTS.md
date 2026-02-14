# Accounts Domain BDD

## Scope
Frontend session authentication and password-change gate behavior.

## Scenarios
- Given a user is authenticated with a Phoenix session, when the user has `must_change_password=true`, then protected frontend project access is blocked.
- Given a user changes password successfully, when the user retries protected frontend project access, then access is allowed.
- Given a session with active organization context, when requesting data for a different org without switching, then access is denied.
- Given a valid switch-organization action, when subsequent requests are made, then authorization uses the new org context.
- Given an invalid or expired session cookie, when calling protected endpoints, then authentication fails.

## Invariants
- Keep existing Phoenix browser auth routes and LiveView auth flow intact.
- Session auth uses secure HTTP-only cookies managed by Phoenix session middleware.
- Password changes must clear `must_change_password`.
