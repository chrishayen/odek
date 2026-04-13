# Requirement: "a date and time utility library for dates, date ranges, periods, and time-of-day"

Calendar math on top of a clock primitive.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
    std.time.unix_to_components
      @ (seconds: i64) -> date_components
      + splits a unix timestamp into year, month, day, hour, minute, second
      # time
    std.time.components_to_unix
      @ (comps: date_components) -> i64
      + reconstructs a unix timestamp from components
      # time

datelib
  datelib.date_new
    @ (year: i32, month: i32, day: i32) -> result[date_value, string]
    + constructs a date value
    - returns error when the day is out of range for the month
    # construction
  datelib.today
    @ () -> date_value
    + returns today's date in local civil terms
    # construction
    -> std.time.now_seconds
    -> std.time.unix_to_components
  datelib.add_days
    @ (date: date_value, days: i32) -> date_value
    + shifts a date by a number of days, crossing months and years
    # arithmetic
    -> std.time.components_to_unix
    -> std.time.unix_to_components
  datelib.days_between
    @ (start: date_value, end: date_value) -> i64
    + returns the signed number of days between two dates
    # arithmetic
  datelib.range_new
    @ (start: date_value, end: date_value) -> result[date_range, string]
    + creates an inclusive date range
    - returns error when end precedes start
    # ranges
  datelib.range_contains
    @ (range: date_range, date: date_value) -> bool
    + reports whether a date falls within the range
    # ranges
  datelib.time_of_day_new
    @ (hour: i32, minute: i32, second: i32) -> result[time_of_day, string]
    + constructs a wall-clock time
    - returns error when any component is out of range
    # time_of_day
  datelib.period_new
    @ (years: i32, months: i32, days: i32) -> period_value
    + describes a calendar period for adding to dates
    # periods
  datelib.add_period
    @ (date: date_value, period: period_value) -> date_value
    + applies a calendar period to a date
    ? month overflow clamps to the last day of the target month
    # arithmetic
