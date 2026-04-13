# Requirement: "a library for attaching hooks to any SQL driver"

Wraps an arbitrary SQL driver with before/after callbacks for queries.

std: (all units exist)

sqlhooks
  sqlhooks.wrap
    @ (driver: sql_driver, hooks: hook_set) -> sql_driver
    + returns a new driver that delegates to the inner one
    + before/after hooks fire for each query execution
    # wrapping
  sqlhooks.hooks_new
    @ () -> hook_set
    + creates an empty hook set
    # construction
  sqlhooks.on_before
    @ (hooks: hook_set, fn: before_fn) -> hook_set
    + registers a callback invoked with the query and args prior to execution
    # registration
  sqlhooks.on_after
    @ (hooks: hook_set, fn: after_fn) -> hook_set
    + registers a callback invoked with the query, args, elapsed time, and error
    + after hooks run even when the query fails
    # registration
