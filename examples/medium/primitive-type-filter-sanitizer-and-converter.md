# Requirement: "a library that filters, sanitizes, and converts values between common primitive types"

Offers type coercion, whitespace/HTML sanitization, and a rule-based filter pipeline over values.

std: (all units exist)

data_filter
  data_filter.to_int
    @ (raw: string) -> result[i64, string]
    + parses a decimal integer, allowing optional leading sign and surrounding whitespace
    - returns error when the input is not a valid integer
    # conversion
  data_filter.to_float
    @ (raw: string) -> result[f64, string]
    + parses a floating-point number, allowing optional sign and exponent
    - returns error on non-numeric input
    # conversion
  data_filter.to_bool
    @ (raw: string) -> result[bool, string]
    + accepts "true"/"false"/"1"/"0"/"yes"/"no" case-insensitively
    - returns error on any other input
    # conversion
  data_filter.trim_spaces
    @ (raw: string) -> string
    + returns the string with leading and trailing ASCII whitespace removed
    # sanitization
  data_filter.strip_tags
    @ (raw: string) -> string
    + removes angle-bracket tags while preserving inner text
    # sanitization
  data_filter.escape_html
    @ (raw: string) -> string
    + replaces &, <, >, ", and ' with their HTML entities
    # sanitization
  data_filter.apply_rules
    @ (value: string, rules: list[string]) -> result[string, string]
    + applies a pipeline of named rules (e.g. "trim|strip_tags|escape_html") left to right
    - returns error when a rule name is unknown
    # pipeline
