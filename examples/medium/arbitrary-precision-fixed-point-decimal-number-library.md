# Requirement: "arbitrary-precision fixed-point decimal numbers"

A decimal is an integer mantissa paired with an exponent. Arithmetic returns new decimals without loss of precision; division takes a requested scale.

std: (all units exist)

decimal
  decimal.from_string
    fn (raw: string) -> result[decimal_value, string]
    + parses "-0.00", "123.456", "1e-5" into an exact decimal
    - returns error on malformed input
    # parsing
  decimal.from_int
    fn (value: i64) -> decimal_value
    + returns a decimal with the given integer value and scale zero
    # construction
  decimal.to_string
    fn (value: decimal_value) -> string
    + renders the decimal in plain (non-scientific) form preserving trailing zeros from scale
    # rendering
  decimal.add
    fn (a: decimal_value, b: decimal_value) -> decimal_value
    + returns the exact sum, rescaling the operand with smaller scale
    # arithmetic
  decimal.sub
    fn (a: decimal_value, b: decimal_value) -> decimal_value
    + returns the exact difference
    # arithmetic
  decimal.mul
    fn (a: decimal_value, b: decimal_value) -> decimal_value
    + returns the exact product with scale equal to the sum of operand scales
    # arithmetic
  decimal.div
    fn (a: decimal_value, b: decimal_value, scale: i32, mode: string) -> result[decimal_value, string]
    + returns the quotient rounded to the requested scale using the named rounding mode
    - returns error on division by zero
    - returns error on unknown rounding mode
    # arithmetic
  decimal.compare
    fn (a: decimal_value, b: decimal_value) -> i32
    + returns -1, 0, or 1 according to numeric ordering regardless of scale
    # comparison
  decimal.rescale
    fn (value: decimal_value, scale: i32, mode: string) -> result[decimal_value, string]
    + returns the value at the requested scale, rounding if the new scale is smaller
    - returns error on unknown rounding mode
    # rescale
