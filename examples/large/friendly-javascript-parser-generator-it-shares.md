# Requirement: "a parser generator that takes a context-free grammar and produces an LALR parser"

Grammars are declared as productions; the generator computes first sets, follow sets, and an LALR(1) parse table; the runtime drives the table over a token stream.

std: (all units exist)

parser_generator
  parser_generator.new_grammar
    @ (start: string) -> grammar_state
    + returns an empty grammar with the given start nonterminal
    # construction
  parser_generator.add_production
    @ (grammar: grammar_state, lhs: string, rhs: list[string]) -> grammar_state
    + appends a production mapping lhs to the ordered rhs symbols
    # grammar
  parser_generator.add_token
    @ (grammar: grammar_state, name: string, pattern: string) -> grammar_state
    + registers a terminal with a matching pattern
    # grammar
  parser_generator.first_set
    @ (grammar: grammar_state, symbol: string) -> list[string]
    + returns the set of terminals that can start derivations of the symbol
    + returns a set containing the empty string when the symbol is nullable
    # analysis
  parser_generator.follow_set
    @ (grammar: grammar_state, nonterminal: string) -> list[string]
    + returns the set of terminals that can immediately follow the nonterminal
    # analysis
  parser_generator.build_lr_items
    @ (grammar: grammar_state) -> list[lr_item_set]
    + returns the canonical collection of LR(1) item sets for the grammar
    # generation
  parser_generator.build_table
    @ (grammar: grammar_state) -> result[parse_table, string]
    + returns the LALR(1) action and goto table
    - returns error when the grammar contains a shift-reduce or reduce-reduce conflict
    # generation
  parser_generator.tokenize
    @ (grammar: grammar_state, source: string) -> result[list[token], string]
    + returns the token stream for the source using registered terminal patterns
    - returns error at the first unmatched character
    # runtime
  parser_generator.parse
    @ (table: parse_table, tokens: list[token]) -> result[parse_node, string]
    + returns a parse tree for a well-formed token stream
    - returns error on an unexpected token with its source position
    # runtime
  parser_generator.generate_runtime_source
    @ (table: parse_table) -> string
    + returns standalone source code that embeds the table and drives parsing
    # codegen
