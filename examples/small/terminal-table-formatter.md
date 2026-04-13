# Requirement: "a tabular data formatter for terminal output"

Lay out rows of cells into an aligned, width-aware text block.

std
  std.text
    std.text.display_width
      @ (s: string) -> i32
      + returns the visible column width of the string
      ? wide glyphs count as two columns
      # text
    std.text.pad_right
      @ (s: string, width: i32) -> string
      + pads the string with spaces on the right up to width
      # text

table
  table.new
    @ (max_width: i32) -> table_state
    + creates an empty table with a maximum total width
    # construction
  table.add_row
    @ (t: table_state, cells: list[string]) -> table_state
    + appends a row of cell strings
    ? rows may have different column counts; missing cells render blank
    # rows
  table.render
    @ (t: table_state) -> string
    + returns the table as aligned text, one row per line
    + columns are sized to the widest cell, clipped so total width fits
    # rendering
    -> std.text.display_width
    -> std.text.pad_right
