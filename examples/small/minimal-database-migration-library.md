# Requirement: "a minimal database migration library"

Two-phase migrations registered in memory and applied in version order.

std: (all units exist)

migrator
  migrator.register
    fn (version: i64, up: migration_fn, down: migration_fn) -> migrator_state
    + adds a migration to the registry
    - overwrites when a version is already registered
    # registry
  migrator.applied_set
    fn (executor: sql_executor) -> result[list[i64], string]
    + returns applied versions, creating the bookkeeping table if absent
    - returns error when the executor rejects the bookkeeping query
    # state
  migrator.migrate
    fn (state: migrator_state, executor: sql_executor) -> result[i32, string]
    + applies unapplied migrations in ascending order and returns the count
    - returns error and halts on the first failure
    # apply
  migrator.rollback_last
    fn (state: migrator_state, executor: sql_executor) -> result[i64, string]
    + rolls back the most recently applied migration and returns its version
    - returns error when no migrations are applied
    # rollback
