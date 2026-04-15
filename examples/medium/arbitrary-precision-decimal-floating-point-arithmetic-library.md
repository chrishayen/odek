# Requirement: "an arbitrary-precision decimal floating-point arithmetic library"

Decimals are represented exactly as (coefficient, exponent) without binary rounding. Division truncates at a caller-supplied scale.

std: (all units exist)

decimal
  decimal.from_string
    fn (value: string) -> result[decimal_value, string]
    + parses signed decimals like "-123.456" and "1e-10"
    - returns error on malformed input
    # parsing
  decimal.to_string
    fn (value: decimal_value) -> string
    + renders the decimal in canonical form without trailing zeros
    # formatting
  decimal.add
    fn (a: decimal_value, b: decimal_value) -> decimal_value
    + aligns exponents and sums the coefficients exactly
    # arithmetic
  decimal.sub
    fn (a: decimal_value, b: decimal_value) -> decimal_value
    + aligns exponents and subtracts exactly
    # arithmetic
  decimal.mul
    fn (a: decimal_value, b: decimal_value) -> decimal_value
    + multiplies coefficients and sums exponents exactly
    # arithmetic
  decimal.div
    fn (a: decimal_value, b: decimal_value, scale: i32) -> result[decimal_value, string]
    + truncated division to the requested number of fractional digits
    - returns error when b is zero
    # arithmetic
  decimal.compare
    fn (a: decimal_value, b: decimal_value) -> i32
    + returns -1, 0, or 1
    # comparison
  decimal.round
    fn (value: decimal_value, scale: i32) -> decimal_value
    + rounds half-away-from-zero to the given scale
    # rounding
