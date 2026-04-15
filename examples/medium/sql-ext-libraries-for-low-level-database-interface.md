# Requirement: "a set of extensions on top of a low-level SQL database interface"

Adds named-parameter queries, row-to-struct mapping, and batch insertion helpers on top of a generic SQL driver.

std
  std.sql
    std.sql.exec
      fn (conn: sql_conn, query: string, args: list[sql_value]) -> result[i64, string]
      + executes a statement and returns the number of affected rows
      - returns error on driver or syntax failure
      # sql
    std.sql.query
      fn (conn: sql_conn, query: string, args: list[sql_value]) -> result[row_iter, string]
      + executes a query and returns an iterator of rows
      - returns error on driver or syntax failure
      # sql
    std.sql.next_row
      fn (iter: row_iter) -> result[optional[map[string, sql_value]], string]
      + returns the next row as a column-to-value map, or none at end
      - returns error on read failure
      # sql

sql_ext
  sql_ext.rewrite_named
    fn (query: string, params: map[string, sql_value]) -> result[tuple[string, list[sql_value]], string]
    + replaces ":name" placeholders with positional markers in source order
    + returns the rewritten query and a matching argument list
    - returns error when a referenced name is missing from params
    # parameter_binding
  sql_ext.named_exec
    fn (conn: sql_conn, query: string, params: map[string, sql_value]) -> result[i64, string]
    + rewrites named parameters and executes the statement
    - returns error on binding or execution failure
    # execution
    -> std.sql.exec
  sql_ext.named_query
    fn (conn: sql_conn, query: string, params: map[string, sql_value]) -> result[list[map[string, sql_value]], string]
    + rewrites named parameters, runs the query, and returns all rows
    - returns error on binding or execution failure
    # execution
    -> std.sql.query
    -> std.sql.next_row
  sql_ext.rows_to_records
    fn (rows: list[map[string, sql_value]], field_map: map[string, string]) -> list[map[string, sql_value]]
    + renames row columns according to the field_map, producing caller-facing records
    # mapping
  sql_ext.bulk_insert
    fn (conn: sql_conn, table: string, columns: list[string], rows: list[list[sql_value]]) -> result[i64, string]
    + builds a multi-values INSERT and returns the number of inserted rows
    - returns error when row lengths do not match the column count
    # bulk_operations
    -> std.sql.exec
