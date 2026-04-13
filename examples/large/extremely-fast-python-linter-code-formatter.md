# Requirement: "a fast source code linter and formatter library"

Covers the full pipeline: lex, parse, lint with pluggable rules, and reprint with a canonical style.

std
  std.text
    std.text.split_lines
      @ (source: string) -> list[string]
      + splits source on line breaks, preserving empty trailing lines
      # text

linter
  linter.tokenize
    @ (source: string) -> result[list[token], string]
    + returns the token stream for the source
    - returns error on an unterminated string or invalid character
    # lexing
  linter.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + returns the abstract syntax tree for the token stream
    - returns error with line and column on a syntax error
    # parsing
  linter.new_rule
    @ (name: string, check: fn(ast_node) -> list[diagnostic]) -> rule
    + creates a lint rule with the given name and checker function
    # rules
  linter.register_rule
    @ (state: linter_state, rule: rule) -> linter_state
    + adds a rule to the linter's active set
    # rules
  linter.disable_rule
    @ (state: linter_state, name: string) -> linter_state
    + removes the named rule from the active set; idempotent
    # rules
  linter.lint
    @ (state: linter_state, source: string) -> result[list[diagnostic], string]
    + runs all active rules against the source and returns ordered diagnostics
    - returns error when tokenization or parsing fails
    # linting
  linter.format
    @ (source: string) -> result[string, string]
    + returns the source reprinted in canonical style, preserving semantics
    - returns error when the source does not parse
    # formatting
  linter.fix
    @ (source: string, diagnostics: list[diagnostic]) -> result[string, string]
    + applies auto-fixes for diagnostics that carry a suggested replacement
    - returns error when two suggested fixes overlap
    # autofix
  linter.diff
    @ (before: string, after: string) -> list[hunk]
    + returns the line-level change hunks between before and after
    # diff
    -> std.text.split_lines
  linter.diagnostic_line
    @ (d: diagnostic) -> string
    + returns a human-readable one-line description of the diagnostic
    # reporting
