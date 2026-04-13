# Requirement: "a database migration tool that supports embedding migrations as in-memory assets"

Applies versioned migrations from an in-memory source, tracking which have been run in a metadata table.

std: (all units exist)

migrator
  migrator.embed_source
    @ (files: map[string, string]) -> migration_source
    + wraps a name-to-sql map as a source of migrations
    # source
  migrator.list_migrations
    @ (source: migration_source) -> list[migration]
    + returns migrations sorted by the numeric prefix of their name
    + ignores entries whose names do not start with digits
    # enumeration
  migrator.applied_versions
    @ (executor: sql_executor) -> result[list[i64], string]
    + returns versions already applied, creating the metadata table if missing
    - returns error when the executor rejects the query
    # state
  migrator.apply
    @ (executor: sql_executor, source: migration_source) -> result[i32, string]
    + runs pending migrations in order and returns how many were applied
    - returns error and halts on the first migration that fails
    # migration
  migrator.rollback
    @ (executor: sql_executor, source: migration_source, steps: i32) -> result[i32, string]
    + rolls back the given number of most-recent migrations
    - returns error when steps exceeds applied count
    # rollback
