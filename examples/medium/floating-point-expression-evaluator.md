# Requirement: "a floating-point arithmetic expression evaluator"

Parses and evaluates infix expressions with +, -, *, /, parentheses, and unary minus.

std: (all units exist)

evaluator
  evaluator.tokenize
    @ (expr: string) -> result[list[token], string]
    + recognizes numbers, operators, and parentheses
    + accepts decimal literals and scientific notation
    - returns error on unexpected character
    # lexing
  evaluator.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + builds an AST respecting standard operator precedence and left-associativity
    + handles unary minus
    - returns error on unbalanced parentheses
    - returns error on trailing operator
    # parsing
  evaluator.eval_ast
    @ (node: ast_node) -> result[f64, string]
    + returns the numeric value of the AST
    - returns error on division by zero
    # evaluation
  evaluator.eval
    @ (expr: string) -> result[f64, string]
    + returns the value of the expression
    - returns error on any lexing, parsing, or evaluation failure
    # facade
    -> evaluator.tokenize
    -> evaluator.parse
    -> evaluator.eval_ast
