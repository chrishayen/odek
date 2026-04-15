# Requirement: "a pretty-printer for arbitrary data structures"

Walks a generic value tree and produces indented multi-line output.

std: (all units exist)

pretty_print
  pretty_print.format
    fn (value: data_value, indent: i32) -> string
    + returns a multi-line, indented textual representation of value
    + handles nested maps and lists recursively
    + renders strings quoted and escaped
    ? indent is the number of spaces per level
    # formatting
  pretty_print.format_compact
    fn (value: data_value) -> string
    + returns a single-line representation of value
    # formatting
  pretty_print.format_diff
    fn (a: data_value, b: data_value) -> string
    + returns a side-by-side textual diff marking added, removed, and changed fields
    + returns empty string when a and b are equal
    # diffing
