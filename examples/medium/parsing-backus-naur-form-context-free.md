# Requirement: "a Backus-Naur form context-free grammar parsing library"

Parses BNF grammar text into rules and provides traversal over productions.

std
  std.strings
    std.strings.split_lines
      @ (s: string) -> list[string]
      + splits on LF, tolerating CRLF
      # strings
    std.strings.trim
      @ (s: string) -> string
      + strips leading and trailing ASCII whitespace
      # strings
    std.strings.starts_with
      @ (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings

bnf
  bnf.parse
    @ (source: string) -> result[bnf_grammar, string]
    + parses lines of "<name> ::= <expression>" into rules
    + joins continuation lines starting with '|'
    - returns error on missing '::=' or unmatched '<' '>'
    # parsing
    -> std.strings.split_lines
    -> std.strings.trim
    -> std.strings.starts_with
  bnf.rules
    @ (g: bnf_grammar) -> list[bnf_rule]
    + returns every parsed rule in source order
    # query
  bnf.rule_for
    @ (g: bnf_grammar, name: string) -> optional[bnf_rule]
    + returns the rule with the given non-terminal name
    # query
  bnf.alternatives
    @ (r: bnf_rule) -> list[list[bnf_symbol]]
    + returns each alternative as a sequence of terminal/non-terminal symbols
    # query
  bnf.non_terminals
    @ (g: bnf_grammar) -> list[string]
    + returns every non-terminal name referenced by the grammar
    # analysis
  bnf.undefined_non_terminals
    @ (g: bnf_grammar) -> list[string]
    + returns non-terminals referenced on the right that have no rule
    # analysis
