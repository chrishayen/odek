# Requirement: "format a date with timezone offset as ISO 8601"

Formats a unix timestamp and offset as an ISO-8601 string like "2015-11-30T10:40:35+01:00".

std
  std.time
    std.time.unix_to_components
      @ (unix_seconds: i64, offset_seconds: i32) -> time_components
      + returns year, month, day, hour, minute, second in the offset's local time
      # time

tzfmt
  tzfmt.format_offset
    @ (offset_seconds: i32) -> string
    + returns the offset as "+HH:MM" or "-HH:MM"
    + returns "+00:00" for offset 0
    - returns "-05:30" for offset -19800
    # formatting
  tzfmt.format
    @ (unix_seconds: i64, offset_seconds: i32) -> string
    + returns an ISO-8601 string like "2015-11-30T10:40:35+01:00"
    + zero-pads every component to its natural width
    # formatting
    -> std.time.unix_to_components
