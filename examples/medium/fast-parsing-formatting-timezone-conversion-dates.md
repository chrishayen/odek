# Requirement: "a library for parsing, formatting, and timezone conversion for dates"

Dates are represented as civil fields plus a zone identifier; the std layer supplies timezone rule lookup and the epoch clock.

std
  std.time
    std.time.load_zone
      @ (name: string) -> result[zone_rules, string]
      + loads IANA zone rules by name (e.g. "America/New_York")
      - returns error when the zone is unknown
      # time
    std.time.unix_to_civil
      @ (unix_seconds: i64, zone: zone_rules) -> civil_time
      + converts a unix timestamp to year/month/day/hour/minute/second in the given zone
      # time
    std.time.civil_to_unix
      @ (civil: civil_time, zone: zone_rules) -> result[i64, string]
      + converts civil fields in a zone to a unix timestamp
      - returns error when the civil time falls in a DST gap
      # time

datetime
  datetime.parse
    @ (text: string, layout: string) -> result[civil_time, string]
    + parses a date string according to the given layout (e.g. "YYYY-MM-DD HH:mm:ss")
    - returns error when the text does not match the layout
    # parsing
  datetime.format
    @ (civil: civil_time, layout: string) -> string
    + renders a civil time using the given layout
    + pads numeric fields to the layout's width
    # formatting
  datetime.convert_zone
    @ (civil: civil_time, from_zone: string, to_zone: string) -> result[civil_time, string]
    + converts a civil time from one zone to another
    - returns error when either zone name is unknown
    # timezone_conversion
    -> std.time.load_zone
    -> std.time.civil_to_unix
    -> std.time.unix_to_civil
