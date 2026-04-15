# Requirement: "a parsing expression grammar parser generator with packrat memoization"

Compiles a PEG grammar into a parser object that can parse inputs with linear-time packrat memoization.

std: (all units exist)

peg
  peg.parse_grammar
    fn (grammar: string) -> result[grammar_state, string]
    + returns a grammar_state for well-formed PEG source
    - returns error with line and column for syntax errors
    ? rule names must be identifiers and may not shadow built-in alternatives
    # grammar_parsing
  peg.compile
    fn (grammar: grammar_state, start_rule: string) -> result[parser_state, string]
    + returns a parser whose entry point is the named rule
    - returns error when the start rule is not defined in the grammar
    - returns error when a referenced rule is never defined
    # compilation
  peg.parse
    fn (parser: parser_state, input: string) -> result[parse_tree, parse_error]
    + returns a parse tree when the input matches the start rule
    - returns a parse_error with position and expected alternatives on failure
    ? uses packrat memoization so each (rule, position) is evaluated at most once
    # parsing
  peg.tree_text
    fn (tree: parse_tree) -> string
    + returns the substring of the original input covered by the tree
    # tree_access
  peg.tree_children
    fn (tree: parse_tree) -> list[parse_tree]
    + returns the child subtrees produced by named rules
    # tree_access
