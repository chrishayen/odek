# Requirement: "an in-process job scheduler with a human-friendly API"

Users register callables to fire on recurring intervals; a tick function advances the scheduler based on the current clock.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

scheduler
  scheduler.new
    @ () -> scheduler_state
    + creates an empty scheduler
    # construction
  scheduler.every_seconds
    @ (state: scheduler_state, seconds: i64, tag: string, job: job_fn) -> scheduler_state
    + registers a job to run every N seconds
    ? first run is scheduled for now + N seconds
    # registration
  scheduler.every_minutes
    @ (state: scheduler_state, minutes: i64, tag: string, job: job_fn) -> scheduler_state
    + registers a job to run every N minutes
    # registration
  scheduler.every_days_at
    @ (state: scheduler_state, hour: i32, minute: i32, tag: string, job: job_fn) -> scheduler_state
    + registers a job to run daily at the given wall-clock time
    - panics when hour or minute are out of range
    # registration
  scheduler.cancel
    @ (state: scheduler_state, tag: string) -> scheduler_state
    + returns a state with the job bearing the given tag removed
    - no-op when no job with that tag exists
    # registration
  scheduler.tick
    @ (state: scheduler_state) -> scheduler_state
    + runs every job whose scheduled time is at or before now, then reschedules
    ? jobs fire in order of their scheduled time
    # execution
    -> std.time.now_seconds
  scheduler.next_run_seconds
    @ (state: scheduler_state) -> optional[i64]
    + returns the unix timestamp of the next scheduled run, or none when empty
    # introspection
