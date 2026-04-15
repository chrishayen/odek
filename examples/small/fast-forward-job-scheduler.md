# Requirement: "a job scheduler with the ability to fast-forward time"

The scheduler holds its own clock so tests can advance virtual time arbitrarily. No real clock dependency.

std: (all units exist)

scheduler
  scheduler.new
    fn (start: i64) -> scheduler_state
    + creates a scheduler with virtual time set to start
    # construction
  scheduler.schedule
    fn (state: scheduler_state, fire_at: i64, action: job_action) -> scheduler_state
    + enqueues action to fire at virtual time fire_at
    # registration
  scheduler.advance_to
    fn (state: scheduler_state, target: i64) -> tuple[list[job_action], scheduler_state]
    + moves virtual time forward to target and returns all actions whose fire_at is <= target, in fire order
    - returns an empty list when target is not greater than current time
    # time_travel
  scheduler.current_time
    fn (state: scheduler_state) -> i64
    + returns the current virtual time
    # introspection
  scheduler.pending
    fn (state: scheduler_state) -> i32
    + returns how many actions remain scheduled
    # introspection
