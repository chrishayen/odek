# Requirement: "a cron-style time range expression parser that checks whether a given time falls within any range"

A range expression is a cron expression paired with a duration; the library answers "is t covered by any range?"

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s on every occurrence of sep
      # strings
    std.strings.trim
      fn (s: string) -> string
      + strips leading and trailing ASCII whitespace
      # strings
  std.time
    std.time.components_utc
      fn (unix_seconds: i64) -> time_components
      + returns minute, hour, day-of-month, month, day-of-week for a UTC instant
      # time

cronrange
  cronrange.parse
    fn (expr: string) -> result[cron_range, string]
    + parses "<cron expression> <duration>" where duration is like "15m" or "2h"
    - returns error when either half is malformed
    # parsing
    -> std.strings.split
    -> std.strings.trim
  cronrange.parse_list
    fn (exprs: list[string]) -> result[list[cron_range], string]
    + parses each expression, returning all ranges or the first error
    - returns error with the index of the first invalid expression
    # parsing
  cronrange.matches
    fn (cr: cron_range, at: i64) -> bool
    + returns true when at falls inside any occurrence of cr
    - returns false when at sits in a gap between occurrences
    # evaluation
    -> std.time.components_utc
  cronrange.any_matches
    fn (ranges: list[cron_range], at: i64) -> bool
    + returns true when at falls inside at least one of the ranges
    # evaluation
  cronrange.next_start
    fn (cr: cron_range, after: i64) -> optional[i64]
    + returns the next instant at which cr begins after the given time
    - returns none when no future occurrence exists within a sensible bound
    # evaluation
