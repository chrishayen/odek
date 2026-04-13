# Requirement: "a parsing expression grammar (PEG) parser generator"

Users build a grammar from primitive combinators, compile it to a parser, and run it over input. Packrat memoization handles left-biased choice without exponential blowup.

std
  std.strings
    std.strings.byte_at
      @ (s: string, i: i32) -> u8
      + returns the byte at index i
      # strings
    std.strings.slice
      @ (s: string, start: i32, end: i32) -> string
      + returns the substring between byte offsets
      # strings
  std.collections
    std.collections.map_get
      @ (m: map[i64, parse_memo], key: i64) -> optional[parse_memo]
      + returns the memoized result for a (rule, position) key
      # collections
    std.collections.map_set
      @ (m: map[i64, parse_memo], key: i64, value: parse_memo) -> map[i64, parse_memo]
      + returns the map with key mapped to value
      # collections

peg
  peg.grammar_new
    @ () -> peg_grammar
    + returns an empty grammar
    # grammar
  peg.rule_literal
    @ (text: string) -> peg_expr
    + matches the exact text
    # combinator
  peg.rule_char_class
    @ (chars: string) -> peg_expr
    + matches any single character listed in chars, with range shorthand
    # combinator
  peg.rule_any
    @ () -> peg_expr
    + matches any single character except end-of-input
    # combinator
  peg.rule_seq
    @ (parts: list[peg_expr]) -> peg_expr
    + matches each part in order, committing as it goes
    # combinator
  peg.rule_choice
    @ (alts: list[peg_expr]) -> peg_expr
    + tries each alternative left-to-right, returning the first that succeeds
    # combinator
  peg.rule_optional
    @ (inner: peg_expr) -> peg_expr
    + matches inner zero or one times
    # combinator
  peg.rule_zero_or_more
    @ (inner: peg_expr) -> peg_expr
    + matches inner as many times as possible
    # combinator
  peg.rule_one_or_more
    @ (inner: peg_expr) -> peg_expr
    + matches inner at least once, then greedily
    # combinator
  peg.rule_and_predicate
    @ (inner: peg_expr) -> peg_expr
    + succeeds without consuming input when inner would succeed
    # combinator
  peg.rule_not_predicate
    @ (inner: peg_expr) -> peg_expr
    + succeeds without consuming input when inner would fail
    # combinator
  peg.rule_ref
    @ (name: string) -> peg_expr
    + references another named rule, resolved at compile time
    # combinator
  peg.grammar_define
    @ (g: peg_grammar, name: string, expr: peg_expr) -> peg_grammar
    + registers a named rule in the grammar
    # grammar
  peg.compile
    @ (g: peg_grammar, start: string) -> result[peg_parser, string]
    + resolves references and returns a runnable parser rooted at start
    - returns error on unresolved rule references
    # compilation
  peg.parse
    @ (p: peg_parser, source: string) -> result[parse_tree, parse_error]
    + runs the parser with packrat memoization, returning a parse tree
    - returns parse_error with position and expected-set on failure
    # execution
    -> std.strings.byte_at
    -> std.strings.slice
    -> std.collections.map_get
    -> std.collections.map_set
  peg.tree_capture
    @ (t: parse_tree, name: string) -> optional[string]
    + returns the substring captured under a named rule
    # tree
  peg.parse_from_dsl
    @ (grammar_source: string) -> result[peg_grammar, string]
    + parses a textual PEG grammar DSL into a peg_grammar
    - returns error with line/column on syntax errors in the DSL
    # parsing
