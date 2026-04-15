# Requirement: "a scanner and parser for LL(1) grammars"

A table-driven LL(1) parser with a hand-rolled lexer. std provides a small token stream utility.

std
  std.text
    std.text.is_alpha
      fn (c: u8) -> bool
      + returns true for ascii letters
      # text
    std.text.is_digit
      fn (c: u8) -> bool
      + returns true for ascii digits 0-9
      # text
    std.text.is_whitespace
      fn (c: u8) -> bool
      + returns true for spaces, tabs, newlines, and carriage returns
      # text

ll1
  ll1.tokenize
    fn (source: string) -> result[list[token], string]
    + produces tokens with kind, lexeme, and byte offset
    - returns error at an unexpected character
    # lexing
    -> std.text.is_alpha
    -> std.text.is_digit
    -> std.text.is_whitespace
  ll1.parse_grammar
    fn (text: string) -> result[grammar, string]
    + reads a BNF grammar description into productions and non-terminal sets
    - returns error on undefined non-terminals
    # grammar
  ll1.first_sets
    fn (g: grammar) -> map[string, list[string]]
    + computes FIRST sets for every non-terminal in the grammar
    # analysis
  ll1.follow_sets
    fn (g: grammar, first: map[string, list[string]]) -> map[string, list[string]]
    + computes FOLLOW sets given the grammar and FIRST sets
    # analysis
  ll1.build_table
    fn (g: grammar) -> result[parse_table, string]
    + constructs the LL(1) parse table
    - returns error when the grammar has a conflict between two productions for one cell
    # table_construction
  ll1.parse
    fn (table: parse_table, tokens: list[token]) -> result[parse_tree, string]
    + drives the predictive parser and returns the parse tree
    - returns error on an unexpected token with expected set
    # parsing
