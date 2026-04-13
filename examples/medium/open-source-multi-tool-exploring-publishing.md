# Requirement: "a multi-tool for exploring and publishing tabular data"

Opens a tabular data file, exposes query and schema introspection, and renders results in a few formats.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when file is missing
      # io
  std.json
    std.json.encode_rows
      @ (rows: list[map[string, string]]) -> string
      + encodes rows as a JSON array of objects
      # serialization
  std.csv
    std.csv.parse
      @ (raw: string) -> result[list[list[string]], string]
      + parses comma-separated rows, handling quoted fields
      - returns error on unterminated quotes
      # parsing

explorer
  explorer.open_csv
    @ (path: string) -> result[table_handle, string]
    + returns a handle whose first row is treated as the header
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
    -> std.csv.parse
  explorer.schema
    @ (t: table_handle) -> list[string]
    + returns the column names in order
    # introspection
  explorer.row_count
    @ (t: table_handle) -> i64
    + returns the number of data rows (excluding the header)
    # introspection
  explorer.select_columns
    @ (t: table_handle, columns: list[string]) -> result[table_handle, string]
    + returns a new handle projected to the given columns
    - returns error when a requested column does not exist
    # query
  explorer.filter_equals
    @ (t: table_handle, column: string, value: string) -> result[table_handle, string]
    + returns rows where the given column equals the given value
    - returns error when the column does not exist
    # query
  explorer.render_json
    @ (t: table_handle) -> string
    + returns rows as a JSON array of objects
    # rendering
    -> std.json.encode_rows
  explorer.render_table
    @ (t: table_handle) -> string
    + returns a text table with aligned columns
    # rendering
