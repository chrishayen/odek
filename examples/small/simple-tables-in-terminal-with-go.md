# Requirement: "a terminal table formatter"

Formats rows of cells into an aligned monospace table with a header and border.

std: (all units exist)

table
  table.new
    @ (headers: list[string]) -> table_state
    + creates a table with the given column headers
    # construction
  table.add_row
    @ (state: table_state, cells: list[string]) -> result[table_state, string]
    + appends a row
    - returns error when the number of cells does not match the header count
    # rows
  table.render
    @ (state: table_state) -> string
    + returns a single string with columns padded to the widest cell in each column
    + draws a top, header-separator, and bottom border using ascii box characters
    ? cell widths are counted in code points, not bytes
    # rendering
