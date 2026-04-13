# Requirement: "a job scheduler supporting webhooks, cron expressions, and one-shot schedules"

A scheduler that fires jobs based on multiple trigger styles. Execution is pluggable; the library owns timing, trigger parsing, and dispatch bookkeeping.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits s on sep
      # strings
    std.strings.trim
      @ (s: string) -> string
      + trims ASCII whitespace from both ends
      # strings

scheduler
  scheduler.parse_cron
    @ (expr: string) -> result[cron_spec, string]
    + parses a five-field cron expression (minute hour day month weekday)
    - returns error on malformed fields or out-of-range values
    # cron
    -> std.strings.split
    -> std.strings.trim
  scheduler.next_cron_fire
    @ (spec: cron_spec, after: i64) -> i64
    + returns the next unix second at which the spec fires strictly after 'after'
    # cron
  scheduler.new
    @ () -> scheduler_state
    + creates an empty scheduler
    # construction
  scheduler.add_cron_job
    @ (state: scheduler_state, id: string, expr: string, action: job_action) -> result[scheduler_state, string]
    + registers a cron-triggered job under id
    - returns error when expr is invalid
    # registration
    -> std.strings.trim
  scheduler.add_oneshot_job
    @ (state: scheduler_state, id: string, fire_at: i64, action: job_action) -> scheduler_state
    + registers a single-fire job at an absolute unix second
    # registration
  scheduler.add_webhook_job
    @ (state: scheduler_state, id: string, secret: string, action: job_action) -> scheduler_state
    + registers a job that fires only when triggered via the webhook trigger function
    # registration
  scheduler.trigger_webhook
    @ (state: scheduler_state, id: string, secret: string) -> result[list[job_action], string]
    + returns the actions to run when the provided secret matches the registered one
    - returns error when id is unknown or secret does not match
    # triggers
  scheduler.due_jobs
    @ (state: scheduler_state) -> tuple[list[job_action], scheduler_state]
    + returns actions for all jobs whose next fire time has passed and advances their schedules
    # dispatch
    -> std.time.now_seconds
  scheduler.remove_job
    @ (state: scheduler_state, id: string) -> scheduler_state
    + removes a job by id; returns state unchanged when id is unknown
    # registration
  scheduler.list_jobs
    @ (state: scheduler_state) -> list[job_summary]
    + returns id, trigger kind, and next fire time for every registered job
    # introspection
