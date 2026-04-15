# Requirement: "a value type conversion library"

Converts between primitive types with overflow and parse checks.

std: (all units exist)

convert
  convert.to_i64
    fn (value: string) -> result[i64, string]
    + parses integer strings like "42" and "-7"
    - returns error for non-numeric input
    - returns error when value overflows i64
    # numeric_parse
  convert.to_f64
    fn (value: string) -> result[f64, string]
    + parses numeric strings into double-precision floats
    - returns error for malformed input like "1.2.3"
    # numeric_parse
  convert.to_bool
    fn (value: string) -> result[bool, string]
    + accepts "true"/"false"/"1"/"0" case-insensitively
    - returns error on unrecognized strings
    # bool_parse
  convert.to_string
    fn (value: f64) -> string
    + renders numeric value using the shortest round-trip representation
    # stringify
