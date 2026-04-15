# Requirement: "a thin query helper layered over a relational database driver"

Makes the common cases (one row, many rows, exec, named parameters) short. The driver primitives live in std.

std
  std.sql
    std.sql.open
      fn (driver: string, dsn: string) -> result[sql_conn, string]
      + opens a database connection
      - returns error on invalid DSN
      # database
    std.sql.exec_raw
      fn (conn: sql_conn, query: string, args: list[sql_value]) -> result[i64, string]
      + executes a statement and returns the number of affected rows
      - returns error on driver failure
      # database
    std.sql.query_raw
      fn (conn: sql_conn, query: string, args: list[sql_value]) -> result[list[map[string, sql_value]], string]
      + runs a query and returns rows as maps keyed by column name
      - returns error on driver failure
      # database

query_helper
  query_helper.rewrite_named_params
    fn (query: string, params: map[string, sql_value]) -> tuple[string, list[sql_value]]
    + replaces :name placeholders with positional markers and returns the args in order
    + preserves colons inside quoted string literals
    # parameters
  query_helper.exec
    fn (conn: sql_conn, query: string, params: map[string, sql_value]) -> result[i64, string]
    + runs a named-parameter statement and returns the affected row count
    - returns error when a placeholder has no matching parameter
    # exec
    -> std.sql.exec_raw
  query_helper.select_one
    fn (conn: sql_conn, query: string, params: map[string, sql_value]) -> result[optional[map[string, sql_value]], string]
    + returns the first row or none if no rows match
    - returns error when the query returns more than one row
    # select_one
    -> std.sql.query_raw
  query_helper.select_all
    fn (conn: sql_conn, query: string, params: map[string, sql_value]) -> result[list[map[string, sql_value]], string]
    + returns all rows as a list of column maps
    # select_all
    -> std.sql.query_raw
  query_helper.in_transaction
    fn (conn: sql_conn, body: tx_callback) -> result[void, string]
    + begins, runs the callback, and commits; rolls back on callback error
    - returns the callback's error after rollback
    # transaction
    -> std.sql.exec_raw
