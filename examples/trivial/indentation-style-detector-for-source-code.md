# Requirement: "a library for detecting the indentation style of source code"

Scans leading whitespace to decide whether a file uses tabs or spaces and the unit width.

std: (all units exist)

detect_indent
  detect_indent.analyze
    @ (source: string) -> indent_style
    + returns the most common leading-whitespace pattern across non-empty lines
    + reports "tab" when tabs dominate, "space" with a width when spaces dominate
    - returns "unknown" when no line has leading whitespace
    # detection
