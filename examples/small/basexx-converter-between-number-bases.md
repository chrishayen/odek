# Requirement: "convert to, from, and between digit strings in various number bases"

Two project functions: parse a digit string in a source base to an integer, and render an integer as a digit string in a target base. Conversion between bases is the composition of the two.

std: (all units exist)

basexx
  basexx.decode
    @ (digits: string, base: i32) -> result[i64, string]
    + decodes "ff" in base 16 to 255
    + decodes "1010" in base 2 to 10
    + accepts both lowercase and uppercase letters for bases above 10
    - returns error when base is outside [2, 36]
    - returns error when a digit is outside the alphabet for the base
    # parsing
  basexx.encode
    @ (value: i64, base: i32) -> result[string, string]
    + encodes 255 in base 16 as "ff"
    + encodes 0 as "0" regardless of base
    + encodes negative values with a leading minus sign
    - returns error when base is outside [2, 36]
    # formatting
