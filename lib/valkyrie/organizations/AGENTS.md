# Organizations Domain BDD

## Scope
Tenancy boundaries and role-based authorization policy.

## Scenarios
- Given organization roles, when listing role definitions, then exactly owner/admin/member/viewer are returned.
- Given cross-organization data access attempts, when requesting resources outside caller organization, then access is denied.
- Given role permission grants, when performing write operations, then authorized roles succeed and unauthorized roles receive forbidden.
- Given membership role changes, when the user performs a subsequent action, then authorization reflects the updated role.

## Invariants
- Core resources are organization scoped.
- Every authorization decision is evaluated against caller role in active organization.
