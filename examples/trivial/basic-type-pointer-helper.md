# Requirement: "helpers for creating pointers to basic-type values"

A tiny helper that wraps a value into an optional so callers can pass it where an optional is expected.

std: (all units exist)

ptrhelp
  ptrhelp.some
    fn (value: i64) -> optional[i64]
    + returns an optional holding the given value
    ? generic over basic types; i64 shown as the canonical case
    # optional_construction
