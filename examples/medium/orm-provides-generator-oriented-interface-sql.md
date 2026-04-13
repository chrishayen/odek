# Requirement: "an object-relational mapper with a generator-oriented query interface"

Queries are built as lazy query streams so iteration can pull rows on demand; the library translates the stream into parameterized SQL at execution time.

std
  std.sql
    std.sql.query
      @ (conn: sql_conn, sql: string, args: list[bytes]) -> result[row_cursor, string]
      + returns a cursor over rows matching the query
      - returns error when the database rejects the sql
      # database
    std.sql.next_row
      @ (cursor: row_cursor) -> result[optional[map[string, bytes]], string]
      + returns the next row as a column-value map or empty when exhausted
      - returns error on driver failure
      # database

gen_orm
  gen_orm.from_table
    @ (table: string) -> query_plan
    + starts a new query plan against a table
    # planning
  gen_orm.filter
    @ (plan: query_plan, column: string, op: string, value: bytes) -> query_plan
    + appends a filter predicate to the plan
    # planning
  gen_orm.select_columns
    @ (plan: query_plan, columns: list[string]) -> query_plan
    + restricts the projection to the given columns
    # planning
  gen_orm.compile
    @ (plan: query_plan) -> tuple[string, list[bytes]]
    + compiles the plan into a parameterized SQL statement and argument list
    # compilation
  gen_orm.iterate
    @ (conn: sql_conn, plan: query_plan) -> result[row_generator, string]
    + executes the plan and returns a generator over result rows
    - returns error when compilation or execution fails
    # iteration
    -> std.sql.query
  gen_orm.next
    @ (gen: row_generator) -> result[optional[map[string, bytes]], string]
    + pulls the next row from the generator
    - returns empty when the generator is exhausted
    # iteration
    -> std.sql.next_row
