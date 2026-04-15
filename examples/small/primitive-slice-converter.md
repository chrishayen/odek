# Requirement: "slice conversion between primitive types"

Element-wise conversions between numeric and string lists. Each function handles one direction and is explicit about overflow behavior.

std: (all units exist)

sliceconv
  sliceconv.ints_to_strings
    fn (xs: list[i64]) -> list[string]
    + returns the decimal string form of each element
    + returns an empty list for an empty input
    # conversion
  sliceconv.strings_to_ints
    fn (xs: list[string]) -> result[list[i64], string]
    + parses each element as a signed decimal integer
    - returns error naming the first element that fails to parse
    # conversion
  sliceconv.floats_to_ints
    fn (xs: list[f64]) -> list[i64]
    + truncates each element toward zero
    # conversion
  sliceconv.ints_to_floats
    fn (xs: list[i64]) -> list[f64]
    + returns each element converted to f64
    # conversion
  sliceconv.bools_to_strings
    fn (xs: list[bool]) -> list[string]
    + returns "true" or "false" for each element
    # conversion
