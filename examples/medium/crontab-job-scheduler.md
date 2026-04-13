# Requirement: "a simple crontab-like job scheduler"

Schedules named jobs using a cron-style expression and runs due jobs when the caller ticks the scheduler.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
    std.time.parse_cron
      @ (expr: string) -> result[cron_spec, string]
      + parses a 5-field cron expression into a spec
      - returns error on invalid fields or wrong number of fields
      # time

scheduler
  scheduler.new
    @ () -> scheduler_state
    + returns an empty scheduler with no jobs
    # construction
  scheduler.add_job
    @ (state: scheduler_state, name: string, cron_expr: string) -> result[scheduler_state, string]
    + adds a job under the given name
    - returns error on duplicate name
    - returns error on invalid cron expression
    # registration
    -> std.time.parse_cron
  scheduler.remove_job
    @ (state: scheduler_state, name: string) -> scheduler_state
    + removes the named job if present
    + is a no-op when the name is unknown
    # registration
  scheduler.next_fire_time
    @ (state: scheduler_state, name: string, after: i64) -> optional[i64]
    + returns the next unix second at which the job would fire
    - returns none when the name is unknown
    # scheduling
  scheduler.tick
    @ (state: scheduler_state) -> tuple[list[string], scheduler_state]
    + returns the names of jobs that are now due and records their fire times
    + never returns the same fire time twice for the same job
    # dispatch
    -> std.time.now_seconds
