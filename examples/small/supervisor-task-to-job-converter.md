# Requirement: "a library that turns a short-lived task into a long-running supervised job"

Wraps a user-provided task and restarts it on failure with backoff. Running the task itself is the caller's job; this layer owns the supervision loop state.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

supervisor
  supervisor.new
    fn (task: fn() -> result[void, string], max_backoff_ms: i64) -> supervisor_state
    + creates a supervisor wrapping task with capped exponential backoff
    # construction
  supervisor.tick
    fn (state: supervisor_state) -> supervisor_state
    + runs one supervision step: invokes task if due, updates backoff on failure, resets on success
    # supervision
    -> std.time.now_millis
  supervisor.status
    fn (state: supervisor_state) -> job_status
    + returns current phase (idle, running, backoff), last error, and restart count
    # introspection
