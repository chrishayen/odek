# Workspace Domain BDD

## Scope
Projects, features, stories, cursor pagination, and soft-delete behavior.

## Scenarios
- Given project definition of done markdown, when reading the project, then the markdown is returned unchanged.
- Given project names in mixed order, when listing projects, then results are sorted alphabetically by name.
- Given an organization already has an active project name (case-insensitive), when another project with the same name is created, then creation fails with uniqueness validation.
- Given feature-linked and featureless stories, when creating stories, then `feature_id` is optional and stored accordingly.
- Given organization-scoped resources, when cross-org reads occur, then access is denied.
- Given paginated story listing, when requesting subsequent pages with cursor, then results continue without duplicates.
- Given new inserts between page requests, when reusing the original cursor, then continuation remains valid.
- Given a soft-deleted story, when listing/reading default endpoints, then it is excluded/not found.
- Given soft-deleted stories with prior activity, when history is requested through allowed access, then historical records remain available.

## Invariants
- Canonical story states: backlog, ready, in_progress, review, done.
- Core entities use soft delete semantics via `deleted_at`.
