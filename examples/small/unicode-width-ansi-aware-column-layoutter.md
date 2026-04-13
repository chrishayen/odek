# Requirement: "a library that lays out a list of strings into columns, aware of unicode width and ansi escape sequences"

Computes the column layout and renders, ignoring ANSI escapes for width purposes.

std: (all units exist)

columns
  columns.visible_width
    @ (s: string) -> i32
    + returns the printed column width, skipping ANSI escape sequences
    + treats wide codepoints as width 2
    # measurement
  columns.plan
    @ (items: list[string], total_width: i32, gap: i32) -> tuple[i32, i32]
    + returns (num_columns, column_width) that fit the widest item
    ? falls back to one column when even a single item cannot fit
    # layout
    -> columns.visible_width
  columns.render
    @ (items: list[string], total_width: i32, gap: i32) -> string
    + returns a newline-separated grid sorted top-to-bottom then left-to-right
    # rendering
    -> columns.plan
    -> columns.visible_width
