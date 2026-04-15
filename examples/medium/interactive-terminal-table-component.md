# Requirement: "an interactive table component for a terminal UI"

Pure state machine for a keyboard-driven table: columns, rows, cursor, selection, sorting. Drawing and input capture belong to the host.

std: (all units exist)

table
  table.new
    fn (columns: list[column_def]) -> table_state
    + creates an empty table with the given column definitions
    # construction
  table.set_rows
    fn (state: table_state, rows: list[list[string]]) -> table_state
    + replaces all rows; each row must have one cell per column
    - rows with mismatched arity are rejected and state is unchanged
    # data
  table.move_cursor
    fn (state: table_state, delta: i32) -> table_state
    + moves the cursor by delta rows, clamped to [0, row_count-1]
    # navigation
  table.toggle_selection
    fn (state: table_state) -> table_state
    + toggles whether the row under the cursor is selected
    # selection
  table.sort_by
    fn (state: table_state, column_index: i32, ascending: bool) -> table_state
    + sorts rows lexicographically by the given column
    - returns state unchanged when column_index is out of range
    # sorting
  table.visible_window
    fn (state: table_state, viewport_rows: i32) -> list[list[string]]
    + returns the slice of rows that should be drawn given the current cursor
    ? window scrolls to keep the cursor in view
    # rendering_query
  table.selected_rows
    fn (state: table_state) -> list[list[string]]
    + returns every row currently marked selected
    # query
