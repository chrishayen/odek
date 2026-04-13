# Requirement: "a parser generator library"

Compiles a context-free grammar into a table-driven Earley-style parser. Users define rules, then feed tokens to get parse trees.

std
  std.collections
    std.collections.set_add
      @ (items: list[string], item: string) -> list[string]
      + returns a new list with item added if not present
      + leaves the list unchanged when item is already present
      # collections

parsergen
  parsergen.new_grammar
    @ (start_symbol: string) -> grammar_state
    + returns an empty grammar with the given start nonterminal
    # construction
  parsergen.add_rule
    @ (grammar: grammar_state, lhs: string, rhs: list[string]) -> grammar_state
    + adds a production lhs -> rhs
    + permits multiple rules with the same lhs
    ? terminals and nonterminals are distinguished by case of the first character
    # grammar
  parsergen.compile
    @ (grammar: grammar_state) -> result[parser, string]
    + returns a parser ready to consume tokens
    - returns error when the start symbol has no rule
    - returns error when any rule references an undefined symbol
    # compilation
  parsergen.feed
    @ (parser: parser, token: string) -> result[parser, string]
    + advances the parser with one token
    - returns error when the token does not match any active expectation
    # parsing
  parsergen.finish
    @ (parser: parser) -> result[list[parse_tree], string]
    + returns all complete parse trees
    - returns error when no complete parse covers the input
    # parsing
  parsergen.parse_all
    @ (grammar: grammar_state, tokens: list[string]) -> result[list[parse_tree], string]
    + compiles the grammar and parses the entire token stream
    - returns error on any compilation or parse failure
    # facade
    -> parsergen.compile
    -> parsergen.feed
    -> parsergen.finish
  parsergen.tree_to_string
    @ (tree: parse_tree) -> string
    + renders a parse tree as an indented s-expression
    # debugging
