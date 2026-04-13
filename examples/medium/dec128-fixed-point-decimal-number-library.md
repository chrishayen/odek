# Requirement: "a 128-bit fixed-point decimal number library"

Fixed-point decimal values stored as a 128-bit mantissa plus a small scale exponent. Arithmetic preserves exact decimal semantics within the representable range.

std
  std.math
    std.math.mul_u64_wide
      @ (a: u64, b: u64) -> tuple[u64, u64]
      + returns (high, low) 128-bit product of two 64-bit unsigned integers
      # arithmetic
    std.math.div_u128_by_u64
      @ (hi: u64, lo: u64, divisor: u64) -> tuple[u64, u64]
      + returns (quotient_lo, remainder) for a 128-by-64 division
      - returns (0, 0) when divisor is 0
      # arithmetic

dec128
  dec128.from_string
    @ (text: string) -> result[dec128_value, string]
    + parses "123.456" into a value with mantissa 123456 and scale 3
    + accepts a leading minus sign
    - returns error on non-digit characters outside the sign and decimal point
    - returns error when the mantissa would overflow 128 bits
    # parsing
  dec128.to_string
    @ (value: dec128_value) -> string
    + renders the value as a decimal string with scale-many fraction digits
    + trims trailing zeros past the decimal point
    # formatting
  dec128.add
    @ (a: dec128_value, b: dec128_value) -> result[dec128_value, string]
    + returns the exact sum, rescaling to the larger scale
    - returns error on 128-bit overflow
    # arithmetic
  dec128.sub
    @ (a: dec128_value, b: dec128_value) -> result[dec128_value, string]
    + returns the exact difference at the larger scale
    - returns error on 128-bit underflow or overflow
    # arithmetic
  dec128.mul
    @ (a: dec128_value, b: dec128_value) -> result[dec128_value, string]
    + returns the product with scale equal to the sum of input scales
    - returns error on 128-bit overflow
    # arithmetic
    -> std.math.mul_u64_wide
  dec128.div
    @ (a: dec128_value, b: dec128_value, result_scale: u8) -> result[dec128_value, string]
    + returns the quotient rounded half-to-even at the requested scale
    - returns error when b is zero
    # arithmetic
    -> std.math.div_u128_by_u64
  dec128.compare
    @ (a: dec128_value, b: dec128_value) -> i8
    + returns -1, 0, or 1 after normalizing to a common scale
    # comparison
