# Requirement: "a high-precision fixed-point decimal library for financial applications"

Fixed-point decimal values optimized for financial arithmetic: exact scale preservation, rounding control, and explicit overflow errors.

std
  std.math
    std.math.round_half_even
      @ (numerator: i128, denominator: i128) -> i128
      + returns numerator/denominator rounded half-to-even
      - returns 0 when denominator is 0
      # arithmetic

udec
  udec.from_string
    @ (text: string) -> result[udec_value, string]
    + parses signed decimal strings like "-123.4567" into exact values
    - returns error on non-numeric characters outside sign and decimal point
    - returns error on mantissa overflow
    # parsing
  udec.from_i64
    @ (whole: i64, scale: u8) -> udec_value
    + constructs a value representing whole * 10^-scale
    # construction
  udec.to_string
    @ (value: udec_value) -> string
    + renders the value preserving its stored scale
    # formatting
  udec.add
    @ (a: udec_value, b: udec_value) -> result[udec_value, string]
    + returns the exact sum at max(scale_a, scale_b)
    - returns error on overflow
    # arithmetic
  udec.sub
    @ (a: udec_value, b: udec_value) -> result[udec_value, string]
    + returns the exact difference at max(scale_a, scale_b)
    - returns error on underflow or overflow
    # arithmetic
  udec.mul
    @ (a: udec_value, b: udec_value) -> result[udec_value, string]
    + returns the product at scale_a + scale_b
    - returns error on overflow
    # arithmetic
  udec.div
    @ (a: udec_value, b: udec_value, result_scale: u8) -> result[udec_value, string]
    + returns the quotient rounded half-to-even at the requested scale
    - returns error on division by zero
    # arithmetic
    -> std.math.round_half_even
  udec.round
    @ (value: udec_value, scale: u8) -> udec_value
    + returns the value rounded half-to-even to the requested scale
    # rounding
    -> std.math.round_half_even
  udec.compare
    @ (a: udec_value, b: udec_value) -> i8
    + returns -1, 0, or 1 after normalizing to a common scale
    # comparison
