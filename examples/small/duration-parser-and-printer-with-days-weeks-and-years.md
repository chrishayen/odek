# Requirement: "a duration library supporting parsing and printing durations in days, weeks, and years"

Durations are represented as a signed number of seconds.

std: (all units exist)

duration
  duration.parse
    @ (text: string) -> result[i64, string]
    + returns total seconds for inputs like "2y3w4d5h6m7s"
    + accepts a leading minus sign for negative durations
    - returns error on unknown unit suffix
    - returns error on empty input
    ? a year is 365 days; a week is 7 days
    # parsing
  duration.format
    @ (seconds: i64) -> string
    + returns a compact form using years, weeks, days, hours, minutes, seconds
    + uses only non-zero units, largest first
    + returns "0s" for zero
    # formatting
  duration.add
    @ (a: i64, b: i64) -> i64
    + returns the sum of two durations in seconds
    # arithmetic
