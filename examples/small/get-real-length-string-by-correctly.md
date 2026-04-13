# Requirement: "a library that returns the visible length of a string, counting astral symbols as one and ignoring ANSI escape codes"

std: (all units exist)

string_length
  string_length.strip_ansi
    @ (s: string) -> string
    + removes CSI-style ANSI escape sequences (ESC [ ... letter)
    + leaves text without escapes unchanged
    # ansi
  string_length.count_code_points
    @ (s: string) -> i32
    + returns the number of unicode code points in s
    ? counts astral (non-BMP) characters as one
    # counting
  string_length.length
    @ (s: string) -> i32
    + returns the visible length: strips ANSI escapes, then counts code points
    # entry
    -> string_length.strip_ansi
    -> string_length.count_code_points
