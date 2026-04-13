# Requirement: "a CSV reader and writer"

Reads records from a CSV string and writes records back out. Quoting and escaping are handled by the reader and writer.

std: (all units exist)

csv
  csv.read_all
    @ (source: string, delimiter: string) -> result[list[list[string]], string]
    + parses the entire input into rows of fields
    - returns error on an unterminated quoted field
    + handles CRLF and LF line endings
    # reading
  csv.write_all
    @ (rows: list[list[string]], delimiter: string) -> string
    + returns a CSV string, quoting fields that contain the delimiter, quotes, or newlines
    # writing
  csv.read_row
    @ (source: string, offset: i32, delimiter: string) -> result[tuple[list[string], i32], string]
    + parses a single row starting at offset and returns the fields and new offset
    - returns error when the row is malformed
    # incremental
  csv.quote_field
    @ (field: string, delimiter: string) -> string
    + returns the field wrapped in quotes when quoting is required, escaping inner quotes
    # quoting
