# Requirement: "a query language that compiles model filter expressions to parameterized SQL"

Parses a small filter language and emits a WHERE clause plus bound parameters.

std: (all units exist)

mql
  mql.tokenize
    @ (source: string) -> result[list[mql_token], string]
    + emits identifiers, string literals, number literals, operators (=, !=, <, <=, >, >=, in, and, or, not), and parentheses
    - returns error on unterminated strings or unknown characters
    # lexing
  mql.parse
    @ (tokens: list[mql_token]) -> result[mql_expr, string]
    + builds an expression tree with standard precedence: not > comparison > and > or
    - returns error on unbalanced parentheses or missing operands
    # parsing
  mql.validate
    @ (expr: mql_expr, allowed_fields: list[string]) -> result[void, string]
    + confirms every identifier in expr appears in allowed_fields
    - returns error naming the first forbidden field
    # validation
  mql.compile
    @ (expr: mql_expr) -> tuple[string, list[value]]
    + renders the expression to a parameterized WHERE fragment and the list of positional bind values
    + uses "$1", "$2", ... as placeholders
    # compilation
  mql.compile_source
    @ (source: string, allowed_fields: list[string]) -> result[tuple[string, list[value]], string]
    + one-shot entry point: tokenize, parse, validate, compile
    - returns error at the first failing stage
    # top_level
