# Requirement: "a database-agnostic query builder with data retrieval helpers"

The library builds SQL text from a fluent query structure and decodes rows into maps. Low-level row iteration is a std primitive.

std
  std.db
    std.db.exec
      @ (dsn: string, sql: string, args: list[string]) -> result[i64, string]
      + executes a statement and returns affected row count
      - returns error on driver or syntax failure
      # database
    std.db.query_rows
      @ (dsn: string, sql: string, args: list[string]) -> result[list[map[string, string]], string]
      + runs a query and returns rows as column-to-value maps
      - returns error on driver or syntax failure
      # database

dbx
  dbx.select
    @ (table: string, columns: list[string]) -> query_state
    + starts a SELECT query with the given columns
    # query_construction
  dbx.where_eq
    @ (q: query_state, column: string, value: string) -> query_state
    + adds an equality predicate bound as a parameter
    # query_construction
  dbx.order_by
    @ (q: query_state, column: string, descending: bool) -> query_state
    + adds an ORDER BY clause
    # query_construction
  dbx.limit
    @ (q: query_state, n: i32) -> query_state
    + sets a row limit
    # query_construction
  dbx.to_sql
    @ (q: query_state) -> tuple[string, list[string]]
    + returns the generated SQL text and the ordered bound parameters
    ? parameters are emitted as positional placeholders so any driver can consume them
    # sql_rendering
  dbx.fetch_all
    @ (dsn: string, q: query_state) -> result[list[map[string, string]], string]
    + executes the query and returns all matching rows
    - returns error when the underlying driver fails
    # retrieval
    -> std.db.query_rows
