# Requirement: "a database management system that stores each table as a text file of line-delimited JSON"

A tiny embedded store where each table is a file of one JSON object per line. Records have integer ids.

std
  std.fs
    std.fs.read_lines
      fn (path: string) -> result[list[string], string]
      + returns file contents split on newlines
      - returns error when file does not exist
      # filesystem
    std.fs.append_line
      fn (path: string, line: string) -> result[void, string]
      + appends a line followed by a newline to the file
      + creates the file when missing
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes string contents to path
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

jsondb
  jsondb.open
    fn (root_dir: string) -> db_state
    + returns a handle rooted at the given directory
    # construction
  jsondb.create_table
    fn (state: db_state, table: string) -> result[void, string]
    + creates an empty table file if it does not exist
    - returns error when the table file exists
    # schema
    -> std.fs.write_all
  jsondb.insert
    fn (state: db_state, table: string, record: map[string, string]) -> result[i64, string]
    + appends the record as a JSON line and returns its assigned id
    - returns error when the table does not exist
    ? ids are monotonic integers starting at 1
    # write
    -> std.json.encode_object
    -> std.fs.append_line
  jsondb.find_by_id
    fn (state: db_state, table: string, id: i64) -> result[optional[map[string, string]], string]
    + returns the record when present
    - returns none when no record has that id
    - returns error when the table does not exist
    # read
    -> std.fs.read_lines
    -> std.json.parse_object
  jsondb.scan
    fn (state: db_state, table: string) -> result[list[map[string, string]], string]
    + returns all records in insertion order
    - returns error when the table does not exist
    # read
    -> std.fs.read_lines
    -> std.json.parse_object
  jsondb.delete_by_id
    fn (state: db_state, table: string, id: i64) -> result[bool, string]
    + returns true when a record was removed
    - returns false when no record has that id
    # write
    -> std.fs.read_lines
    -> std.fs.write_all
