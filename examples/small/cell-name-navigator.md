# Requirement: "a spreadsheet cell name navigator"

Converts between spreadsheet cell names (like "A1", "AB12") and zero-based (row, column) coordinates, and offsets a cell by a delta.

std: (all units exist)

cell_nav
  cell_nav.parse_cell_name
    @ (name: string) -> result[tuple[i32, i32], string]
    + returns (row, column) for a cell name like "B3"
    + handles multi-letter columns like "AA" and "ZZ"
    - returns error on empty string
    - returns error on lowercase or non-alphanumeric input
    # parsing
  cell_nav.format_cell_name
    @ (row: i32, col: i32) -> result[string, string]
    + returns the cell name for a given row and column
    - returns error when row or col is negative
    # formatting
  cell_nav.offset
    @ (name: string, row_delta: i32, col_delta: i32) -> result[string, string]
    + returns the cell name offset by the given deltas
    - returns error when the resulting coordinates are negative
    # navigation
  cell_nav.column_name_to_index
    @ (letters: string) -> result[i32, string]
    + converts a column letter sequence to a zero-based index
    - returns error on non-letter input
    # parsing
  cell_nav.column_index_to_name
    @ (index: i32) -> result[string, string]
    + converts a zero-based column index to its letter sequence
    - returns error when index is negative
    # formatting
