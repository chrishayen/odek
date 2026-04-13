# Requirement: "a type conversion utility library"

Converts between common primitive types with predictable failure semantics. No reflection tricks exposed to the caller.

std: (all units exist)

conv
  conv.to_i64
    @ (value: string) -> result[i64, string]
    + parses a signed decimal integer within i64 range
    - returns error on empty string or non-digit characters
    - returns error when the parsed value overflows i64
    # conversion
  conv.to_f64
    @ (value: string) -> result[f64, string]
    + parses a decimal floating point number
    - returns error on malformed input
    # conversion
  conv.to_bool
    @ (value: string) -> result[bool, string]
    + recognizes "true"/"false"/"1"/"0" case-insensitively
    - returns error on any other input
    # conversion
  conv.i64_to_string
    @ (value: i64) -> string
    + renders a signed integer as its canonical decimal form
    # conversion
  conv.f64_to_string
    @ (value: f64) -> string
    + renders a float without trailing zeros
    ? uses the shortest round-trip representation
    # conversion
