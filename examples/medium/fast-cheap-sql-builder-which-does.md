# Requirement: "a SQL query builder"

Builds SELECT/INSERT/UPDATE/DELETE strings with parameter placeholders. No execution, no escaping of data — parameters are returned separately.

std: (all units exist)

sql_builder
  sql_builder.select
    @ (table: string, columns: list[string]) -> query_state
    + starts a SELECT over the given columns and table
    + uses "*" when columns is empty
    # select
  sql_builder.where_eq
    @ (state: query_state, column: string, value: string) -> query_state
    + appends an equality predicate and a positional parameter
    ? multiple calls are AND-ed together
    # where
  sql_builder.order_by
    @ (state: query_state, column: string, ascending: bool) -> query_state
    + appends an ORDER BY clause
    # ordering
  sql_builder.limit
    @ (state: query_state, n: i64) -> query_state
    + sets the LIMIT
    # limit
  sql_builder.insert
    @ (table: string, values: map[string, string]) -> query_state
    + builds an INSERT with parameter placeholders for each value
    - returns an empty state when values is empty
    # insert
  sql_builder.update
    @ (table: string, assignments: map[string, string]) -> query_state
    + builds an UPDATE with placeholders for each assignment
    # update
  sql_builder.render
    @ (state: query_state) -> tuple[string, list[string]]
    + returns the final SQL text and the parameters in order
    # render
