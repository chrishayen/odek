# Requirement: "a SQL query builder with first-class null handling"

A fluent builder that emits parameterized SQL. Null values in inputs become SQL NULL rather than an empty string or zero.

std: (all units exist)

sql_builder
  sql_builder.select
    fn (table: string) -> query_state
    + starts a SELECT * FROM table query
    # construction
  sql_builder.columns
    fn (q: query_state, cols: list[string]) -> query_state
    + replaces the projection with the given columns
    # projection
  sql_builder.where_eq
    fn (q: query_state, column: string, value: optional[sql_value]) -> query_state
    + adds column = ? with a bound parameter when value is present
    + adds column IS NULL with no parameter when value is absent
    # filtering
  sql_builder.where_in
    fn (q: query_state, column: string, values: list[optional[sql_value]]) -> query_state
    + splits the list into present and absent, emitting column IN (?, ?, ...) OR column IS NULL as needed
    - produces 1 = 0 when the list is empty, so no rows match
    # filtering
  sql_builder.order_by
    fn (q: query_state, column: string, descending: bool) -> query_state
    + appends an ORDER BY clause in the given direction
    # ordering
  sql_builder.build
    fn (q: query_state) -> tuple[string, list[sql_value]]
    + returns the final SQL string and the ordered bound parameters
    + parameter list excludes nulls that were folded into IS NULL
    # rendering
