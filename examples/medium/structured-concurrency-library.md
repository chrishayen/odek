# Requirement: "a structured concurrency library"

Spawn child tasks inside a scope; the scope does not return until all children are done, and any child error cancels siblings.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

structconc
  structconc.new_scope
    @ () -> scope_state
    + creates a new task scope with no children
    # construction
  structconc.spawn
    @ (scope: scope_state, task_id: string, work: task_fn) -> scope_state
    + schedules work under the scope, tagged with task_id
    ? task_fn is an opaque callable that takes a cancel_token and returns result[void, string]
    # spawn
  structconc.wait_all
    @ (scope: scope_state) -> result[map[string, string], string]
    + blocks until all children complete successfully, returning their results
    - returns error holding the first failure and cancels siblings
    # join
  structconc.cancel_scope
    @ (scope: scope_state) -> scope_state
    + signals cancellation to every running child
    # cancel
  structconc.scope_deadline
    @ (scope: scope_state, deadline_millis: i64) -> scope_state
    + attaches a wall-clock deadline; triggers cancel when exceeded
    # deadline
    -> std.time.now_millis
  structconc.is_cancelled
    @ (token: cancel_token) -> bool
    + tasks poll this to cooperatively exit early
    # cancel
