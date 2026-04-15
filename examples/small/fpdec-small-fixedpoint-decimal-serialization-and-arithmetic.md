# Requirement: "serialization and arithmetic for small fixed-point decimals"

Values are stored as a scaled integer with an implicit fixed scale; arithmetic preserves the scale.

std: (all units exist)

fpdec
  fpdec.from_string
    fn (text: string, scale: i32) -> result[fpdec_value, string]
    + parses a decimal string like "3.14" at the given scale
    - returns error on non-numeric input
    - returns error when more fractional digits are supplied than the scale allows
    # parsing
  fpdec.to_string
    fn (value: fpdec_value) -> string
    + renders the value with its fixed number of fractional digits
    # formatting
  fpdec.add
    fn (a: fpdec_value, b: fpdec_value) -> result[fpdec_value, string]
    + sums two values at the same scale
    - returns error when scales differ
    # arithmetic
  fpdec.sub
    fn (a: fpdec_value, b: fpdec_value) -> result[fpdec_value, string]
    + subtracts b from a at the same scale
    - returns error when scales differ
    # arithmetic
  fpdec.mul_int
    fn (value: fpdec_value, factor: i64) -> fpdec_value
    + multiplies the value by an integer factor, preserving scale
    # arithmetic
