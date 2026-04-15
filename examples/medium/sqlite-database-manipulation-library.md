# Requirement: "a library for manipulating sqlite databases"

Wraps the underlying sqlite driver with ergonomic table and row operations. Raw query execution is a thin std primitive.

std
  std.sqlite
    std.sqlite.open
      fn (path: string) -> result[sqlite_conn, string]
      + opens (or creates) a sqlite database file
      - returns error when the file cannot be created or is corrupt
      # database
    std.sqlite.exec
      fn (conn: sqlite_conn, sql: string, params: list[string]) -> result[void, string]
      + executes a statement with positional parameters
      - returns error on syntax or constraint failure
      # database
    std.sqlite.query
      fn (conn: sqlite_conn, sql: string, params: list[string]) -> result[list[map[string, string]], string]
      + executes a query and returns rows as string-keyed maps
      - returns error on syntax failure
      # database

sqlite_tools
  sqlite_tools.open_db
    fn (path: string) -> result[sqlite_conn, string]
    + opens a database at the given path
    # database_access
    -> std.sqlite.open
  sqlite_tools.create_table
    fn (conn: sqlite_conn, table: string, columns: map[string, string]) -> result[void, string]
    + creates a table mapping column_name to sql type if it does not exist
    - returns error when the table exists with a conflicting schema
    # schema
    -> std.sqlite.exec
  sqlite_tools.insert_row
    fn (conn: sqlite_conn, table: string, row: map[string, string]) -> result[void, string]
    + inserts a row; columns missing from the map are left null
    - returns error when the table does not exist
    # rows
    -> std.sqlite.exec
  sqlite_tools.upsert_row
    fn (conn: sqlite_conn, table: string, pk_column: string, row: map[string, string]) -> result[void, string]
    + inserts or replaces by primary key column
    - returns error when pk_column is not present in row
    # rows
    -> std.sqlite.exec
  sqlite_tools.find_rows
    fn (conn: sqlite_conn, table: string, filter: map[string, string]) -> result[list[map[string, string]], string]
    + returns all rows where every filter column equals its value
    # rows
    -> std.sqlite.query
  sqlite_tools.delete_rows
    fn (conn: sqlite_conn, table: string, filter: map[string, string]) -> result[i32, string]
    + deletes matching rows and returns the number removed
    # rows
    -> std.sqlite.exec
  sqlite_tools.list_tables
    fn (conn: sqlite_conn) -> result[list[string], string]
    + returns the names of all user tables
    # schema
    -> std.sqlite.query
