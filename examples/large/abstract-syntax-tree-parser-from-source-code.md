# Requirement: "a parser that produces an abstract syntax tree from source code"

A small lexer plus recursive-descent parser for a C-like expression and statement language. Output is an AST the caller can walk.

std
  std.strings
    std.strings.is_digit
      fn (c: u8) -> bool
      + true when c is an ASCII digit
      # strings
    std.strings.is_alpha
      fn (c: u8) -> bool
      + true when c is an ASCII letter or underscore
      # strings

parser
  parser.tokenize
    fn (source: string) -> result[list[token], parse_error]
    + returns the token stream for a well-formed source string
    - returns error with line and column on an unterminated string literal
    # lexing
    -> std.strings.is_digit
    -> std.strings.is_alpha
  parser.parse_expression
    fn (tokens: list[token]) -> result[expr_node, parse_error]
    + parses an expression with standard precedence for arithmetic and comparison
    - returns error on unexpected token
    # parsing
  parser.parse_statement
    fn (tokens: list[token]) -> result[stmt_node, parse_error]
    + parses assignment, if, while, return, and expression statements
    - returns error when no rule matches the leading token
    # parsing
  parser.parse_program
    fn (source: string) -> result[program_node, parse_error]
    + parses a full source string into an ordered list of top-level statements
    - returns error when trailing tokens remain after the last statement
    # parsing
  parser.walk
    fn (root: program_node, visit: fn(ast_node) -> void) -> void
    + invokes visit on every node in pre-order
    # traversal
  parser.source_range
    fn (n: ast_node) -> tuple[i32, i32]
    + returns (start_offset, end_offset) of a node in the original source
    # diagnostics
