# Requirement: "a library for interacting with intervals of time, defined as a start time and a duration"

A small value-type library. Intervals are (start, duration) pairs; operations check containment and overlap.

std: (all units exist)

timespan
  timespan.new
    @ (start_unix_ms: i64, duration_ms: i64) -> timespan_value
    + constructs an interval with the given start and duration
    - duration_ms must be non-negative; negative values are clamped to zero
    # construction
  timespan.end
    @ (span: timespan_value) -> i64
    + returns start + duration as a unix ms timestamp
    # accessor
  timespan.contains
    @ (span: timespan_value, t_unix_ms: i64) -> bool
    + returns true when t lies within [start, start+duration)
    - returns false for t equal to the exclusive end
    # containment
  timespan.overlaps
    @ (a: timespan_value, b: timespan_value) -> bool
    + returns true when the two half-open intervals share any instant
    - returns false when one ends exactly where the other begins
    # overlap
  timespan.intersection
    @ (a: timespan_value, b: timespan_value) -> optional[timespan_value]
    + returns the overlapping interval when they overlap
    - returns none when there is no overlap
    # intersection
