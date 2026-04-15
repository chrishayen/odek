# Requirement: "a library that expands macro invocations in source code"

Takes source text and a macro table, replaces each invocation with its expansion, and recurses until no macros remain.

std: (all units exist)

macroexp
  macroexp.new_table
    fn () -> macro_table
    + creates an empty macro table
    # construction
  macroexp.define
    fn (table: macro_table, name: string, params: list[string], body: string) -> macro_table
    + registers a macro with named parameters and a body template
    # registration
  macroexp.find_invocation
    fn (source: string, start: i32) -> optional[invocation]
    + returns the next macro invocation starting at or after start, including its byte range and argument list
    - returns none when no invocation is found
    # parsing
  macroexp.substitute
    fn (body: string, params: list[string], args: list[string]) -> string
    + returns the body with each parameter replaced by its argument
    ? parameter references are whole tokens, not substrings
    # substitution
  macroexp.expand_once
    fn (source: string, table: macro_table) -> string
    + replaces every top-level invocation in source with its substituted body
    - leaves text unchanged when no registered macros are invoked
    # expansion
  macroexp.expand_fully
    fn (source: string, table: macro_table, max_depth: i32) -> result[string, string]
    + repeatedly expands until a fixed point is reached
    - returns error when expansion depth exceeds max_depth (likely recursive macro)
    # fixed_point
