# Requirement: "a SQL query builder that composes SELECT statements from chained clauses"

A fluent builder that accumulates clauses and emits a parameterized SQL string and its argument list.

std: (all units exist)

query
  query.select
    fn (table: string, columns: list[string]) -> query_state
    + creates a builder with the given table and projection
    ? empty columns list means "SELECT *"
    # construction
  query.where_eq
    fn (q: query_state, column: string, value: string) -> query_state
    + appends an equality condition, ANDed with existing ones
    # filtering
  query.order_by
    fn (q: query_state, column: string, direction: string) -> result[query_state, string]
    + appends ORDER BY with "asc" or "desc"
    - returns error when direction is not "asc" or "desc"
    # ordering
  query.limit
    fn (q: query_state, n: i64) -> result[query_state, string]
    + sets LIMIT n
    - returns error when n < 0
    # pagination
  query.build
    fn (q: query_state) -> tuple[string, list[string]]
    + returns the final SQL string and the positional parameters
    + placeholders are numbered "$1", "$2", ... in the order conditions were added
    # compilation
