# Requirement: "a library for parsing cron expressions, scheduling tasks, and running a daemon over a crontab-like list"

Parses five-field cron expressions, computes the next fire time, and runs a loop that dispatches tasks as their times elapse.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns the current unix time in seconds
      # time
    std.time.sleep_seconds
      @ (seconds: i64) -> void
      + suspends the caller for the given number of seconds
      # time
    std.time.to_calendar
      @ (unix_seconds: i64) -> calendar_time
      + returns minute, hour, day, month, weekday components in UTC
      # time

cron
  cron.parse
    @ (expr: string) -> result[cron_spec, string]
    + parses a five-field cron expression (minute hour dom month dow)
    + accepts comma lists, ranges with dashes, stepped values with slashes, and star
    - returns error when the expression does not have five fields
    - returns error when any field is out of range
    # parsing
  cron.matches
    @ (spec: cron_spec, at: calendar_time) -> bool
    + true when the calendar time satisfies the spec
    # matching
  cron.next_after
    @ (spec: cron_spec, from_unix: i64) -> i64
    + returns the earliest unix second strictly after from_unix that matches
    # scheduling
    -> std.time.to_calendar
  cron.new_runner
    @ () -> runner_state
    + creates an empty runner with no scheduled tasks
    # construction
  cron.schedule
    @ (state: runner_state, name: string, expr: string, task: task_handler) -> result[runner_state, string]
    + adds a named task with a cron expression and handler
    - returns error when the expression is invalid
    - returns error when a task with that name exists
    # configuration
    -> cron.parse
  cron.tick
    @ (state: runner_state, now: i64) -> runner_state
    + fires every task whose next-fire time is <= now and reschedules it
    # execution
    -> cron.next_after
  cron.run_forever
    @ (state: runner_state) -> result[void, string]
    + loops forever: computes next fire time, sleeps, dispatches, repeats
    - returns error when a task handler signals fatal failure
    # daemon
    -> std.time.now_seconds
    -> std.time.sleep_seconds
