# Requirement: "a fluent SQL query builder with parameterized values"

Queries are constructed immutably: each builder call returns a new query. The final render produces a SQL string with placeholders and a separate argument list.

std: (all units exist)

sqlb
  sqlb.select
    @ (columns: list[string]) -> query
    + starts a SELECT query with the given column list
    ? empty columns render as "*"
    # construction
  sqlb.from
    @ (q: query, table: string) -> query
    + sets the FROM clause
    # from
  sqlb.where_eq
    @ (q: query, column: string, value: sql_value) -> query
    + adds a column = ? predicate with the value captured as an argument
    + multiple where clauses are combined with AND
    # filter
  sqlb.where_in
    @ (q: query, column: string, values: list[sql_value]) -> query
    + adds a column IN (?, ?, ...) predicate
    - produces a clause that is always false when values is empty
    # filter
  sqlb.order_by
    @ (q: query, column: string, ascending: bool) -> query
    + appends an ORDER BY clause
    # ordering
  sqlb.limit
    @ (q: query, n: i64) -> query
    + sets a LIMIT
    # limit
  sqlb.insert
    @ (table: string, columns: list[string], values: list[sql_value]) -> query
    + builds an INSERT with parameterized values
    - produces an error query when columns and values lengths differ
    # construction
  sqlb.render
    @ (q: query) -> result[tuple[string, list[sql_value]], string]
    + returns the SQL text with ? placeholders and the ordered argument list
    - returns error when required clauses like FROM are missing
    # rendering
