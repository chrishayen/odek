# Requirement: "a cron job scheduler with job queuing and live monitoring"

Parses cron expressions, schedules jobs, enqueues runs when workers are busy, and exposes live status snapshots.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
    std.time.sleep_millis
      fn (millis: i64) -> void
      + blocks the current thread for the given duration
      # time
  std.sync
    std.sync.new_mutex
      fn () -> mutex
      + creates an unlocked mutex
      # concurrency
    std.sync.lock
      fn (m: mutex) -> void
      + acquires the mutex, blocking if held
      # concurrency
    std.sync.unlock
      fn (m: mutex) -> void
      + releases a held mutex
      # concurrency

cron_scheduler
  cron_scheduler.parse_expression
    fn (expr: string) -> result[cron_spec, string]
    + parses a five-field cron expression into minute/hour/dom/month/dow sets
    + accepts "*", ranges "a-b", steps "*/n", and comma lists
    - returns error when any field is out of range
    - returns error when the expression does not have exactly five fields
    # parsing
  cron_scheduler.next_fire_after
    fn (spec: cron_spec, after_unix: i64) -> i64
    + returns the next unix timestamp at or after the given time that matches the spec
    ? second precision; always advances at least one minute
    # scheduling
  cron_scheduler.new
    fn (worker_count: i32, queue_capacity: i32) -> scheduler_state
    + creates a scheduler with a fixed worker pool and bounded job queue
    # construction
    -> std.sync.new_mutex
  cron_scheduler.register_job
    fn (state: scheduler_state, name: string, expr: string, handler_id: string) -> result[string, string]
    + registers a job under a unique name and returns its id
    - returns error when the cron expression is invalid
    - returns error when the name is already registered
    # registration
  cron_scheduler.tick
    fn (state: scheduler_state) -> i32
    + fires all jobs whose next_fire_after has passed, enqueueing runs
    + returns the number of runs enqueued this tick
    # scheduling_loop
    -> std.time.now_seconds
  cron_scheduler.dequeue_run
    fn (state: scheduler_state) -> optional[job_run]
    + returns the next pending run for a worker to execute, or none if empty
    # queue
    -> std.sync.lock
    -> std.sync.unlock
  cron_scheduler.record_run_result
    fn (state: scheduler_state, run_id: string, success: bool, duration_millis: i64) -> void
    + updates the history and metrics for a completed run
    # monitoring
  cron_scheduler.snapshot
    fn (state: scheduler_state) -> scheduler_snapshot
    + returns a copy of current job list, queue depth, running count, and last-run summaries
    ? snapshot is read-consistent under the internal mutex
    # monitoring
  cron_scheduler.remove_job
    fn (state: scheduler_state, name: string) -> bool
    + removes a registered job by name; returns true if a job was removed
    # registration
  cron_scheduler.drain
    fn (state: scheduler_state) -> void
    + waits until the queue is empty and all workers are idle
    # lifecycle
    -> std.time.sleep_millis
