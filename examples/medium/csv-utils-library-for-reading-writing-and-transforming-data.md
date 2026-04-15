# Requirement: "a library of utilities for reading, writing, and transforming CSV data"

Reading and writing go through one std primitive each; project-level runes handle typed projection and simple analytics on parsed rows.

std
  std.csv
    std.csv.parse
      fn (text: string) -> result[list[list[string]], string]
      + parses CSV text into rows of string fields, honoring quoted fields and escaped quotes
      - returns error on unterminated quoted fields
      # parsing
    std.csv.serialize
      fn (rows: list[list[string]]) -> string
      + serializes rows into CSV text, quoting fields that contain delimiters or newlines
      # serialization

csv_utils
  csv_utils.from_text
    fn (text: string, has_header: bool) -> result[table, string]
    + parses CSV text into a table; when has_header is true the first row becomes the column names
    - returns error when rows have inconsistent column counts
    # parsing
    -> std.csv.parse
  csv_utils.to_text
    fn (t: table) -> string
    + serializes a table back to CSV text including the header row
    # serialization
    -> std.csv.serialize
  csv_utils.select_columns
    fn (t: table, names: list[string]) -> result[table, string]
    + returns a new table containing only the named columns in the given order
    - returns error when any name is not present
    # projection
  csv_utils.filter_rows
    fn (t: table, pred: fn(map[string, string]) -> bool) -> table
    + returns a table with only the rows for which pred returns true
    # selection
  csv_utils.convert_column
    fn (t: table, column: string, target: string) -> result[table, string]
    + coerces each value in the column to the target type (int, float, bool, date)
    - returns error when any value cannot be parsed
    # typing
  csv_utils.summarize_column
    fn (t: table, column: string) -> result[column_summary, string]
    + returns count, nulls, min, max, and mean for a numeric column
    - returns error when the column is not numeric
    # analytics
  csv_utils.join
    fn (left: table, right: table, on: string) -> result[table, string]
    + inner-joins two tables on the named column
    - returns error when the column is missing in either side
    # joining
