# Requirement: "a source-code AST inspection toolkit"

Lex, parse, and walk a small C-family language into an abstract syntax tree the caller can inspect.

std: (all units exist)

ast
  ast.tokenize
    fn (source: string) -> result[list[token], string]
    + returns the token stream for identifiers, numbers, strings, operators, and keywords
    - returns error on an unterminated string literal
    - returns error on an unrecognized character
    # lexing
  ast.parse_file
    fn (tokens: list[token]) -> result[file_node, string]
    + parses a token stream into a file node containing declarations
    - returns error with position on the first syntactic mismatch
    # parsing
  ast.declarations
    fn (file: file_node) -> list[decl_node]
    + returns the top-level declarations in source order
    # query
  ast.walk
    fn (node: ast_node, visit: fn(ast_node) -> bool) -> void
    + invokes visit on every descendant; stops descending when visit returns false
    # traversal
  ast.find_by_name
    fn (file: file_node, name: string) -> optional[decl_node]
    + returns the first top-level declaration matching name
    - returns none when no declaration matches
    # query
    -> ast.walk
  ast.pretty_print
    fn (node: ast_node) -> string
    + returns a formatted source representation of the node with two-space indentation
    # formatting
  ast.position_of
    fn (node: ast_node) -> source_position
    + returns the 1-based line and column where the node begins
    # query
