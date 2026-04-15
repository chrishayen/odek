# Requirement: "define time-based recurring tasks and compute their next fire times"

A task has a cron-like schedule and a name. The scheduler is pure: given the current time, it returns which tasks are due and updates its internal last-fire bookkeeping.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

scheduler
  scheduler.new
    fn () -> scheduler_state
    + creates an empty scheduler
    # construction
  scheduler.parse_schedule
    fn (expr: string) -> result[schedule, string]
    + parses a five-field cron expression into a schedule
    + supports "*", exact values, ranges "a-b", and steps "*/n"
    - returns error on malformed expression
    - returns error when any field is out of range
    # parsing
  scheduler.register
    fn (state: scheduler_state, name: string, sched: schedule) -> result[scheduler_state, string]
    + adds a task with a schedule
    - returns error when the name is already registered
    # registration
  scheduler.next_fire
    fn (sched: schedule, after: i64) -> i64
    + returns the next unix second at which the schedule matches strictly after the given time
    # scheduling
  scheduler.due
    fn (state: scheduler_state, now: i64) -> tuple[list[string], scheduler_state]
    + returns the names of tasks whose next fire time is at or before now, and the updated state
    ? last-fire times are advanced so a task fires at most once per match
    # scheduling
    -> std.time.now_seconds
  scheduler.remove
    fn (state: scheduler_state, name: string) -> result[scheduler_state, string]
    + removes a registered task
    - returns error when the name is not registered
    # registration
