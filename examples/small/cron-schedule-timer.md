# Requirement: "a ticker that fires at times matching a cron schedule"

Parse a cron expression once, then repeatedly compute the next matching instant. Time primitives live in std.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
    std.time.components
      fn (unix_seconds: i64) -> time_parts
      + returns (minute, hour, day_of_month, month, day_of_week) in UTC
      # time

cron_ticker
  cron_ticker.parse
    fn (expression: string) -> result[cron_spec, string]
    + parses a five-field cron expression with numbers, ranges, lists, and steps
    - returns error when any field is out of range or malformed
    # parsing
  cron_ticker.next_after
    fn (spec: cron_spec, from_unix: i64) -> i64
    + returns the first unix timestamp strictly after from_unix whose components all match
    ? result is aligned to a whole minute
    # scheduling
    -> std.time.components
  cron_ticker.next
    fn (spec: cron_spec) -> i64
    + returns the next firing time after the current clock
    # scheduling
    -> std.time.now_seconds
