# Requirement: "a job task scheduler with worker pool, cron expressions, dates, and human-readable schedule syntax"

A scheduler that accepts multiple schedule description formats and dispatches jobs to a worker pool. Time is injected; worker execution is pluggable.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
    std.time.parse_iso8601
      fn (s: string) -> result[i64, string]
      + parses an ISO 8601 timestamp to unix seconds
      - returns error on malformed input
      # time
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s on sep
      # strings
    std.strings.to_lower
      fn (s: string) -> string
      + returns an ASCII-lowercased copy
      # strings
    std.strings.trim
      fn (s: string) -> string
      + trims ASCII whitespace from both ends
      # strings

scheduler
  scheduler.parse_cron
    fn (expr: string) -> result[schedule, string]
    + parses a five-field cron expression
    - returns error on bad fields
    # parsing
    -> std.strings.split
    -> std.strings.trim
  scheduler.parse_human
    fn (phrase: string) -> result[schedule, string]
    + parses phrases like "every 5 minutes", "every day at 9", "every monday at noon"
    - returns error when the phrase cannot be interpreted
    # parsing
    -> std.strings.to_lower
    -> std.strings.split
  scheduler.parse_date
    fn (iso: string) -> result[schedule, string]
    + treats an ISO 8601 timestamp as a one-shot schedule at that moment
    - returns error on malformed input
    # parsing
    -> std.time.parse_iso8601
  scheduler.next_fire
    fn (s: schedule, after: i64) -> optional[i64]
    + returns the next unix second the schedule fires strictly after 'after'
    - returns none when a one-shot schedule has already fired
    # scheduling
  scheduler.new
    fn (worker_count: i32) -> scheduler_state
    + creates a scheduler with a worker pool of the given size
    # construction
  scheduler.register
    fn (state: scheduler_state, id: string, sched: schedule, task: task_fn) -> scheduler_state
    + registers a task under id with its schedule
    # registration
  scheduler.unregister
    fn (state: scheduler_state, id: string) -> scheduler_state
    + removes a registered task
    # registration
  scheduler.tick
    fn (state: scheduler_state) -> tuple[list[worker_assignment], scheduler_state]
    + advances scheduling: returns (worker_index, task_fn) pairs for jobs whose fire time has passed and workers are free
    # dispatch
    -> std.time.now_seconds
  scheduler.complete
    fn (state: scheduler_state, worker_index: i32) -> scheduler_state
    + marks a worker as free again after its assigned task returns
    # dispatch
  scheduler.pending
    fn (state: scheduler_state) -> i32
    + returns the number of tasks with a scheduled next fire
    # introspection
  scheduler.busy_workers
    fn (state: scheduler_state) -> i32
    + returns how many workers are currently occupied
    # introspection
