# Requirement: "a parser generator that produces parsers with precise error reporting"

Compiles a PEG-style grammar into a parser closure, then applies it to inputs.

std: (all units exist)

peg
  peg.parse_grammar
    @ (source: string) -> result[grammar_ast, string]
    + parses grammar source into an abstract syntax tree
    - returns error with line and column on unexpected tokens
    # grammar_parsing
  peg.compile
    @ (ast: grammar_ast) -> result[parser, string]
    + compiles a grammar ast into a runnable parser
    - returns error when a rule references an undefined name
    - returns error when the start rule is missing
    # compilation
  peg.parse
    @ (p: parser, input: string) -> result[parse_tree, parse_error]
    + returns a parse tree when the input matches the start rule
    - returns a parse_error with line, column, expected tokens, and the furthest-reached position on failure
    # parsing
  peg.format_error
    @ (err: parse_error, input: string) -> string
    + returns a human-readable error with a caret pointing at the failure column
    # diagnostics
  peg.walk
    @ (tree: parse_tree, visitor: tree_visitor) -> void
    + invokes visitor callbacks in depth-first order over the tree
    # traversal
