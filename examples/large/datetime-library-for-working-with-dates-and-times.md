# Requirement: "a library for working with dates and times"

Construction, arithmetic, formatting, parsing, and duration math for civil date-times.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.days_in_month
      fn (year: i32, month: i32) -> i32
      + returns the number of days in the given month
      + accounts for leap years in February
      # calendar
    std.time.is_leap_year
      fn (year: i32) -> bool
      + returns true for Gregorian leap years
      # calendar

datetime
  datetime.new
    fn (year: i32, month: i32, day: i32, hour: i32, minute: i32, second: i32) -> result[datetime, string]
    + returns a datetime when all components are valid
    - returns error when any component is out of range
    # construction
    -> std.time.days_in_month
  datetime.now
    fn () -> datetime
    + returns the current datetime in UTC
    # construction
    -> std.time.now_millis
  datetime.from_unix_millis
    fn (ms: i64) -> datetime
    + constructs a UTC datetime from a unix millisecond timestamp
    # construction
  datetime.to_unix_millis
    fn (dt: datetime) -> i64
    + returns the unix millisecond timestamp
    # conversion
  datetime.add_days
    fn (dt: datetime, days: i32) -> datetime
    + returns dt shifted forward by the given number of days
    + accepts negative values for the past
    # arithmetic
    -> std.time.days_in_month
  datetime.add_seconds
    fn (dt: datetime, seconds: i64) -> datetime
    + returns dt shifted forward by the given number of seconds
    # arithmetic
  datetime.difference_seconds
    fn (a: datetime, b: datetime) -> i64
    + returns a - b in whole seconds
    # arithmetic
  datetime.start_of_day
    fn (dt: datetime) -> datetime
    + returns midnight on dt's date
    # rounding
  datetime.weekday
    fn (dt: datetime) -> i32
    + returns 0 for Sunday through 6 for Saturday
    # calendar
  datetime.format
    fn (dt: datetime, pattern: string) -> result[string, string]
    + formats using tokens like YYYY, MM, DD, HH, mm, ss
    - returns error on unknown tokens
    # formatting
  datetime.parse
    fn (input: string, pattern: string) -> result[datetime, string]
    + parses input according to the same token set
    - returns error when input does not match the pattern
    # parsing
