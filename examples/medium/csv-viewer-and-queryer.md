# Requirement: "a library for viewing and querying CSV files"

Parses a CSV into typed rows, lets callers run simple filter/select expressions, and formats a sliced view for display.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full contents of the file as a string
      - returns error when the file does not exist or cannot be read
      # filesystem

csv_view
  csv_view.parse
    @ (raw: string) -> result[table, string]
    + parses header and rows, handles quoted fields and escaped quotes
    - returns error when a row has a different column count than the header
    # parsing
  csv_view.load
    @ (path: string) -> result[table, string]
    + reads the file and parses it
    - returns error on read or parse failure
    # loading
    -> std.fs.read_all
  csv_view.infer_types
    @ (t: table) -> table
    + promotes columns whose values all parse as integers or floats
    # typing
  csv_view.filter
    @ (t: table, expr: string) -> result[table, string]
    + keeps rows matching a predicate like "col op value" with op in =, !=, <, <=, >, >=
    - returns error on unknown column or malformed expression
    # querying
  csv_view.select
    @ (t: table, columns: list[string]) -> result[table, string]
    + projects the named columns in the given order
    - returns error on unknown column name
    # querying
  csv_view.slice
    @ (t: table, offset: i32, limit: i32) -> table
    + returns a view with up to limit rows starting at offset
    # viewing
  csv_view.render
    @ (t: table) -> string
    + returns an aligned text table with column headers
    # rendering
