# Requirement: "a terminal table rendering library with color support"

Build a table model and render it as a text block. Color codes are plain ANSI escapes the caller can strip if desired.

std: (all units exist)

table
  table.new
    @ (headers: list[string]) -> table_state
    + returns an empty table with the given column headers
    # construction
  table.add_row
    @ (t: table_state, cells: list[string]) -> result[table_state, string]
    + appends the row to the table
    - returns error when the cell count does not match the header count
    # rows
  table.set_column_color
    @ (t: table_state, column: i32, color: string) -> table_state
    + returns a table with the color applied to cells in the column
    # styling
  table.render
    @ (t: table_state) -> string
    + returns a formatted, padded, color-annotated string
    ? columns are padded to the width of the widest cell
    # rendering
