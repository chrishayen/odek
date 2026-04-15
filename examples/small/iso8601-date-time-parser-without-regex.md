# Requirement: "an efficient ISO 8601 date-time parser without regex"

Hand-rolled character walk, no regex engine. Splits into date, time, and offset phases.

std: (all units exist)

iso8601
  iso8601.parse
    fn (input: string) -> result[iso8601_instant, string]
    + parses "2024-03-17T12:34:56Z" into year, month, day, hour, minute, second, nanos, offset_minutes
    + parses fractional seconds with any precision between 1 and 9 digits
    + accepts "+HH:MM", "-HH:MM", and "Z" timezone designators
    - returns error when the date and time separator is neither 'T' nor ' '
    - returns error when month, day, hour, minute, or second fields are out of range
    ? nanos normalized to a 9-digit scale regardless of input precision
    # parsing
  iso8601.format
    fn (instant: iso8601_instant) -> string
    + emits the canonical "YYYY-MM-DDTHH:MM:SS.fffffffffZ" form
    + omits the fractional component when nanos is zero
    # formatting
  iso8601.to_unix_seconds
    fn (instant: iso8601_instant) -> i64
    + converts a parsed instant to seconds since the unix epoch, applying the offset
    # conversion
