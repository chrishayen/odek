# Requirement: "a job scheduler"

Jobs are registered with a delay or interval, then the scheduler is ticked forward in time to run any due jobs. Time comes through a std primitive so tests can drive it deterministically.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

scheduler
  scheduler.new
    @ () -> scheduler_state
    + creates an empty scheduler with no jobs
    # construction
  scheduler.schedule_once
    @ (state: scheduler_state, name: string, run_at_millis: i64) -> result[scheduler_state, string]
    + returns a new state containing the one-shot job
    - returns error when a job with the same name already exists
    # registration
  scheduler.schedule_every
    @ (state: scheduler_state, name: string, interval_millis: i64, first_run_at_millis: i64) -> result[scheduler_state, string]
    + returns a new state containing the recurring job
    - returns error when interval_millis is zero or negative
    # registration
  scheduler.cancel
    @ (state: scheduler_state, name: string) -> scheduler_state
    + removes the named job if present
    + returns unchanged state when the name is unknown
    # registration
  scheduler.due_jobs
    @ (state: scheduler_state, now_millis: i64) -> list[string]
    + returns names of all jobs whose run time is at or before now
    + returns an empty list when nothing is due
    # dispatch
  scheduler.advance
    @ (state: scheduler_state, now_millis: i64) -> scheduler_state
    + marks one-shot jobs that have fired as removed
    + reschedules recurring jobs to their next interval tick
    # dispatch
    -> std.time.now_millis
