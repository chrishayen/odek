# Requirement: "an in-process task scheduler that runs registered functions on cron and interval triggers"

The scheduler owns a store of jobs keyed by id and computes next-fire times from triggers. std supplies monotonic time and cron expression parsing.

std
  std.time
    std.time.now_unix
      @ () -> i64
      + returns the current unix timestamp in seconds
      # time
  std.cron
    std.cron.parse
      @ (expr: string) -> result[cron_spec, string]
      + parses a standard 5-field cron expression
      - returns error on malformed fields or out-of-range values
      # parsing
    std.cron.next_after
      @ (spec: cron_spec, after: i64) -> i64
      + returns the next unix timestamp matching the spec strictly after the given time
      # cron

scheduler
  scheduler.new
    @ () -> scheduler_state
    + creates an empty scheduler with no registered jobs
    # construction
  scheduler.add_cron
    @ (s: scheduler_state, id: string, expr: string) -> result[void, string]
    + registers a job that fires on the given cron expression
    - returns error when id is already registered
    - returns error when expr is invalid
    # registration
    -> std.cron.parse
  scheduler.add_interval
    @ (s: scheduler_state, id: string, period_sec: i64) -> result[void, string]
    + registers a job that fires every period_sec seconds starting now + period
    - returns error when period_sec <= 0
    - returns error when id is already registered
    # registration
  scheduler.remove
    @ (s: scheduler_state, id: string) -> bool
    + removes the job and returns true if it existed
    # registration
  scheduler.due
    @ (s: scheduler_state) -> list[string]
    + returns ids of jobs whose next-fire time is at or before now
    + updates each returned job's next-fire time to the subsequent slot
    # dispatch
    -> std.time.now_unix
    -> std.cron.next_after
  scheduler.next_fire_at
    @ (s: scheduler_state, id: string) -> optional[i64]
    + returns the unix timestamp of the given job's next scheduled fire
    - returns none when the id is unknown
    # inspection
