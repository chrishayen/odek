# Requirement: "an implementation of the Bailey-Borwein-Plouffe algorithm for hexadecimal digits of pi"

BBP lets you extract a hex digit of pi at any position without computing preceding digits. The project exposes one extraction function built from a modular-exponentiation helper and a series term.

std: (all units exist)

pihex
  pihex.mod_pow
    @ (base: i64, exponent: i64, modulus: i64) -> i64
    + returns base^exponent mod modulus via square-and-multiply
    ? used inside the BBP fractional series to keep intermediate values bounded
    # modular_arithmetic
  pihex.series
    @ (j: i32, n: i64) -> f64
    + returns the fractional part of the sum sum_k 16^(n-k) / (8k + j) for k = 0..n plus the tail
    ? j is 1, 4, 5, or 6 per the BBP formula
    # bbp_series
    -> pihex.mod_pow
  pihex.hex_digit_at
    @ (n: i64) -> string
    + returns the single hex digit of pi at position n (0-indexed, after the point)
    ? combines four series values: s1 - s4 - s5 - s6 times appropriate coefficients
    # digit_extraction
    -> pihex.series
  pihex.hex_digits
    @ (start: i64, count: i32) -> string
    + returns a string of count hex digits of pi beginning at position start
    - returns "" when count <= 0
    # digit_sequence
    -> pihex.hex_digit_at
