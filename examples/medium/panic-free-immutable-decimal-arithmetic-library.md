# Requirement: "an immutable decimal number library with panic-free arithmetic"

Decimals as (coefficient, exponent). All ops return result so overflow and division-by-zero are surfaced as errors instead of panics.

std: (all units exist)

decimal
  decimal.from_string
    @ (s: string) -> result[decimal_value, string]
    + parses "123.456" into an exact decimal
    + accepts leading sign and scientific notation
    - returns error on non-numeric input
    # parsing
  decimal.from_int
    @ (n: i64) -> decimal_value
    + returns a decimal with exponent 0 and coefficient n
    # construction
  decimal.to_string
    @ (d: decimal_value) -> string
    + returns the canonical textual form with the stored scale preserved
    # formatting
  decimal.add
    @ (a: decimal_value, b: decimal_value) -> result[decimal_value, string]
    + returns a new decimal; inputs unchanged
    + aligns exponents before adding
    - returns error when the result coefficient overflows i128
    # arithmetic
  decimal.sub
    @ (a: decimal_value, b: decimal_value) -> result[decimal_value, string]
    + returns a - b as a new decimal
    - returns error on coefficient overflow
    # arithmetic
  decimal.mul
    @ (a: decimal_value, b: decimal_value) -> result[decimal_value, string]
    + returns a * b as a new decimal with exponent = ea + eb
    - returns error on coefficient overflow
    # arithmetic
  decimal.div
    @ (a: decimal_value, b: decimal_value, scale: i32) -> result[decimal_value, string]
    + returns a / b rounded half-even to the given scale
    - returns error when b is zero
    - returns error when scale is negative
    # arithmetic
  decimal.compare
    @ (a: decimal_value, b: decimal_value) -> i32
    + returns -1, 0, or 1 after aligning exponents
    # comparison
  decimal.rescale
    @ (d: decimal_value, scale: i32) -> result[decimal_value, string]
    + returns the same value expressed at the target scale
    - returns error when rescaling would require rounding to satisfy the new scale
    # rescale
