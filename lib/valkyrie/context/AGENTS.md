# Context Domain BDD

## Scope
Project context aggregation from persisted workspace and chat activity.

## Scenarios
- Given a project with stories, comments, history, threads, and messages, when requesting project context, then all those records are included.
- Given cross-organization project context requests, when caller org does not own project, then access is denied.
- Given soft-deleted stories, when requesting default context, then deleted stories are excluded.
- Given v1 resource catalog, when listing resources, then dedicated memory/context resource tables are absent.

## Invariants
- Context is derived from persisted core records only.
- No separate memory/context tables in v1.
