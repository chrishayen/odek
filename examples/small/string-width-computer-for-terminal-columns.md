# Requirement: "a library for computing the visual display width of a string in terminal columns"

Classifies each code point as zero-width, narrow, or wide and sums the widths.

std: (all units exist)

string_width
  string_width.code_point_width
    @ (cp: i32) -> i32
    + returns 0 for combining marks and zero-width joiners
    + returns 2 for code points in wide CJK and emoji ranges
    + returns 1 for ordinary printable code points
    - returns 0 for control characters
    # classification
  string_width.measure
    @ (s: string) -> i32
    + returns the total column width by summing per-code-point widths
    + returns 0 for the empty string
    # measurement
