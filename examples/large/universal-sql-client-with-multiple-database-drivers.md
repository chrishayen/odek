# Requirement: "a universal SQL client library supporting multiple database drivers"

Core for a database client library: driver registry, connection via URL, query execution, result formatting, and an interactive REPL-style input handler. No TTY or CLI; the caller drives it.

std
  std.url
    std.url.parse
      fn (raw: string) -> result[url_parts, string]
      + splits a URL into scheme, host, port, path, and query parameters
      - returns error on malformed input
      # parsing
  std.io
    std.io.read_line
      fn (reader: reader_handle) -> result[optional[string], string]
      + reads a single line; returns none at EOF
      # io

usql
  usql.new
    fn () -> client_state
    + creates a client with no registered drivers or open connections
    # construction
  usql.register_driver
    fn (state: client_state, scheme: string, driver: driver) -> client_state
    + registers a driver to handle URLs with the given scheme
    # drivers
  usql.connect
    fn (state: client_state, url: string) -> result[tuple[conn_id, client_state], string]
    + parses the URL, selects the driver by scheme, and opens a connection
    - returns error when no driver is registered for the scheme
    - returns error when the driver rejects the credentials
    # connection
    -> std.url.parse
  usql.disconnect
    fn (state: client_state, conn: conn_id) -> result[client_state, string]
    + closes the connection and removes it from the state
    - returns error when the id is unknown
    # connection
  usql.exec
    fn (state: client_state, conn: conn_id, sql: string) -> result[i64, string]
    + runs a non-returning statement and returns affected row count
    - returns error when the driver reports a syntax or permission failure
    # execution
  usql.query
    fn (state: client_state, conn: conn_id, sql: string) -> result[result_set, string]
    + runs a SELECT and returns columns plus rows
    - returns error on driver failure
    # execution
  usql.format_table
    fn (rs: result_set) -> string
    + renders a result set as an aligned ASCII table with header separator
    + returns "(no rows)" when the set is empty
    # formatting
  usql.format_csv
    fn (rs: result_set) -> string
    + renders a result set as CSV, quoting fields that contain commas or newlines
    # formatting
  usql.split_statements
    fn (script: string) -> list[string]
    + splits a script into individual statements on semicolons outside quoted strings
    # parsing
  usql.read_multiline
    fn (reader: reader_handle) -> result[optional[string], string]
    + reads lines until a statement terminator and returns the assembled text
    + returns none on EOF with no pending content
    # input
    -> std.io.read_line
  usql.describe_table
    fn (state: client_state, conn: conn_id, table: string) -> result[list[column_info], string]
    + returns the column metadata for the named table
    - returns error when the table does not exist
    # introspection
