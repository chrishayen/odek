# Requirement: "a function to reverse a string"

One pure function. No helpers.

std: (all units exist)

string_utils
  string_utils.reverse
    @ (s: string) -> string
    + returns "olleh" when given "hello"
    + returns "" when given ""
    + reverses by grapheme cluster, not by byte, to handle utf-8 correctly
    ? unicode normalization is not performed
    # string_manipulation
