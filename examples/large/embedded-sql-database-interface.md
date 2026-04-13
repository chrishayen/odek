# Requirement: "an embedded SQL database interface"

A cursor/connection API over an embedded engine. The engine itself is a std primitive; the project layer exposes the familiar connection-cursor-row shape.

std
  std.sqlengine
    std.sqlengine.open
      @ (path: string) -> result[engine_handle, string]
      + opens or creates a database file at the given path
      + path ":memory:" creates an in-memory database
      - returns error when the file cannot be created
      # engine
    std.sqlengine.close
      @ (handle: engine_handle) -> result[void, string]
      + releases all resources held by the engine handle
      # engine
    std.sqlengine.prepare
      @ (handle: engine_handle, sql: string) -> result[statement_handle, string]
      + compiles a SQL statement and returns a reusable handle
      - returns error on syntax errors
      # engine
    std.sqlengine.bind
      @ (stmt: statement_handle, index: i32, value: string) -> result[statement_handle, string]
      + binds a positional parameter by 1-based index
      - returns error when index is out of range
      # engine
    std.sqlengine.step
      @ (stmt: statement_handle) -> result[step_outcome, string]
      + advances the statement; outcome is one of "row", "done", or "error"
      # engine
    std.sqlengine.column_string
      @ (stmt: statement_handle, index: i32) -> optional[string]
      + returns the text value of a column by 0-based index
      - returns none when the column is null
      # engine
    std.sqlengine.reset
      @ (stmt: statement_handle) -> statement_handle
      + resets a prepared statement for reuse
      # engine
    std.sqlengine.finalize
      @ (stmt: statement_handle) -> result[void, string]
      + releases resources held by a prepared statement
      # engine

sql
  sql.connect
    @ (path: string) -> result[connection_state, string]
    + opens a connection to a database at the given path
    - returns error when the database cannot be opened
    # connection
    -> std.sqlengine.open
  sql.close
    @ (conn: connection_state) -> result[void, string]
    + closes the connection and all open cursors
    # connection
    -> std.sqlengine.close
  sql.cursor
    @ (conn: connection_state) -> cursor_state
    + creates a new cursor bound to the connection
    # cursor
  sql.execute
    @ (cur: cursor_state, sql: string, params: list[string]) -> result[cursor_state, string]
    + prepares, binds, and steps a statement in one call
    - returns error on SQL syntax or parameter mismatch
    # execution
    -> std.sqlengine.prepare
    -> std.sqlengine.bind
    -> std.sqlengine.step
  sql.execute_many
    @ (cur: cursor_state, sql: string, rows: list[list[string]]) -> result[cursor_state, string]
    + executes the same statement once per parameter row
    # execution
    -> std.sqlengine.prepare
    -> std.sqlengine.bind
    -> std.sqlengine.step
    -> std.sqlengine.reset
  sql.fetch_one
    @ (cur: cursor_state) -> optional[list[string]]
    + returns the next row as a list of column strings
    - returns none when the cursor is exhausted
    # result_fetching
    -> std.sqlengine.step
    -> std.sqlengine.column_string
  sql.fetch_all
    @ (cur: cursor_state) -> list[list[string]]
    + returns every remaining row
    # result_fetching
    -> std.sqlengine.step
    -> std.sqlengine.column_string
  sql.row_count
    @ (cur: cursor_state) -> i32
    + returns the number of rows affected by the last execute
    # result_fetching
  sql.commit
    @ (conn: connection_state) -> result[connection_state, string]
    + commits the active transaction
    # transaction
    -> std.sqlengine.prepare
    -> std.sqlengine.step
  sql.rollback
    @ (conn: connection_state) -> result[connection_state, string]
    + rolls back the active transaction
    # transaction
    -> std.sqlengine.prepare
    -> std.sqlengine.step
