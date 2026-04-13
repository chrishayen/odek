# Requirement: "an extension library for the SQL database package adding named queries, struct scanning, and batch operations"

Thin wrappers over raw SQL that handle named-parameter rewriting, row-to-map scanning, and bulk execution.

std
  std.sql
    std.sql.exec
      @ (conn: sql_conn, sql: string, args: list[string]) -> result[i64, string]
      + executes a statement with positional parameters and returns the affected row count
      - returns error on driver or syntax failure
      # sql
    std.sql.query
      @ (conn: sql_conn, sql: string, args: list[string]) -> result[list[map[string, string]], string]
      + runs a query with positional parameters and returns the rows as string-to-string maps
      - returns error on driver or syntax failure
      # sql

sqlx
  sqlx.rewrite_named
    @ (sql: string, params: map[string, string]) -> result[tuple[string, list[string]], string]
    + replaces :name placeholders with positional markers and returns the argument list in order
    - returns error when a :name placeholder has no matching parameter
    - returns error when a parameter is provided but never referenced
    # named_parameters
  sqlx.exec_named
    @ (conn: sql_conn, sql: string, params: map[string, string]) -> result[i64, string]
    + rewrites named parameters and executes, returning the affected row count
    - propagates errors from rewrite or execution
    # named_parameters
    -> std.sql.exec
  sqlx.query_named
    @ (conn: sql_conn, sql: string, params: map[string, string]) -> result[list[map[string, string]], string]
    + rewrites named parameters and runs the query, returning rows as maps
    - propagates errors from rewrite or execution
    # named_parameters
    -> std.sql.query
  sqlx.scan_one
    @ (rows: list[map[string, string]], key: string) -> result[map[string, string], string]
    + returns the row whose key column matches when there is exactly one match
    - returns error when zero or more than one row matches
    # struct_scanning
  sqlx.batch_exec
    @ (conn: sql_conn, sql: string, rows: list[map[string, string]]) -> result[i64, string]
    + executes the same named statement once per row and returns the total affected count
    - returns error and stops on the first failing row
    # batch
    -> std.sql.exec
