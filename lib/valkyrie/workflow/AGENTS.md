# Workflow Domain BDD

## Scope
Story polling, claim workflow, state transitions, comments, and history events.

## Scenarios
- Given stories eligible for work, when polling, then only eligible stories are returned with cursor paging.
- Given concurrent claim attempts, when two callers claim the same story, then exactly one succeeds.
- Given a claimed story, when reading story fields, then no canonical assignee field is exposed.
- Given valid canonical target states, when updating story state, then any-to-any transitions are accepted.
- Given invalid target states, when updating state, then validation errors are returned.
- Given story field updates, when patching tracked fields, then structured `field_update` history events are recorded.
- Given comment creation, when adding comments, then comments are immutable and `comment_added` history is recorded.

## Invariants
- Comments are append-only.
- History events include type, actor identity, timestamp, and payload metadata.
