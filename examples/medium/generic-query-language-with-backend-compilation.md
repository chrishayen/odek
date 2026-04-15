# Requirement: "a generic query language that compiles to backend-specific query expressions"

The language is parsed into an AST; a compiler walks the AST and produces a target query plus its parameters.

std: (all units exist)

query
  query.parse
    fn (source: string) -> result[query_expr, string]
    + parses source like `age > 18 and name = "ada"` into an expression tree
    - returns error on unbalanced parentheses
    - returns error on unknown operator
    # parsing
  query.validate
    fn (expr: query_expr, schema: map[string, string]) -> result[void, string]
    + checks that every referenced column exists in schema
    - returns error when a column is unknown
    - returns error when an operator is incompatible with the column type
    # validation
  query.compile_sql
    fn (expr: query_expr) -> tuple[string, list[string]]
    + emits a parameterized SQL WHERE fragment and its argument list
    # compilation
  query.compile_filter
    fn (expr: query_expr) -> tuple[string, list[string]]
    + emits a document-store filter expression and its argument list
    # compilation
  query.evaluate
    fn (expr: query_expr, row: map[string, string]) -> result[bool, string]
    + evaluates the expression against an in-memory row
    - returns error when a referenced column is missing from the row
    # evaluation
