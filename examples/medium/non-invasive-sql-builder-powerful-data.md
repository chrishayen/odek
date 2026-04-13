# Requirement: "a sql query builder with parameter binding"

Build SELECT/INSERT/UPDATE/DELETE statements and return the final query string with a parameter list.

std: (all units exist)

sqlb
  sqlb.select
    @ (table: string) -> query_state
    + creates a SELECT query targeting the given table
    # construction
  sqlb.columns
    @ (q: query_state, cols: list[string]) -> query_state
    + sets the column projection for a SELECT query
    + defaults to "*" when called with an empty list
    # projection
  sqlb.where_eq
    @ (q: query_state, column: string, value: sql_value) -> query_state
    + adds an equality predicate joined with AND
    # filtering
  sqlb.where_in
    @ (q: query_state, column: string, values: list[sql_value]) -> query_state
    + adds an IN predicate with one placeholder per value
    - adds a predicate that is always false when values is empty
    # filtering
  sqlb.order_by
    @ (q: query_state, column: string, descending: bool) -> query_state
    + appends an ORDER BY clause
    # ordering
  sqlb.limit
    @ (q: query_state, n: i32) -> query_state
    + sets the LIMIT clause
    # pagination
  sqlb.insert
    @ (table: string, row: map[string, sql_value]) -> query_state
    + creates an INSERT query with the given columns and values
    - returns an empty-valued query when row is empty
    # construction
  sqlb.update
    @ (table: string, set: map[string, sql_value]) -> query_state
    + creates an UPDATE query with SET assignments
    # construction
  sqlb.delete
    @ (table: string) -> query_state
    + creates a DELETE query targeting the given table
    # construction
  sqlb.build
    @ (q: query_state) -> tuple[string, list[sql_value]]
    + returns the final parameterized sql string and the ordered parameter list
    + uses positional placeholders in declaration order
    # emission
