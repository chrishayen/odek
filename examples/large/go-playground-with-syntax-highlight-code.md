# Requirement: "a code playground library with syntax highlighting and code completion"

Edits a source buffer, tokenizes for highlighting, offers completion candidates, and runs snippets through a pluggable executor.

std
  std.text
    std.text.split_lines
      @ (s: string) -> list[string]
      + splits on LF, keeping empty trailing lines
      # text
    std.text.prefix_match
      @ (candidates: list[string], prefix: string) -> list[string]
      + returns candidates that start with prefix
      # text
  std.lex
    std.lex.tokenize
      @ (source: string, grammar: grammar) -> list[token]
      + returns positioned tokens for the given grammar
      + recognizes keywords, identifiers, literals, operators, and comments
      # lexing
    std.lex.default_grammar
      @ () -> grammar
      + returns a grammar covering a C-like syntax family
      # lexing

playground
  playground.new
    @ () -> buffer
    + creates an empty source buffer at cursor 0
    # construction
  playground.insert
    @ (buf: buffer, at: i32, text: string) -> buffer
    + inserts text at the given byte offset
    - returns the buffer unchanged when at is out of range
    # editing
  playground.delete_range
    @ (buf: buffer, start: i32, end: i32) -> buffer
    + removes the byte range [start, end)
    - returns the buffer unchanged when the range is invalid
    # editing
  playground.source
    @ (buf: buffer) -> string
    + returns the current source text
    # inspection
  playground.highlight
    @ (buf: buffer) -> list[highlighted_span]
    + returns spans annotated with token class for every token in the buffer
    # highlighting
    -> std.lex.tokenize
    -> std.lex.default_grammar
  playground.line_column
    @ (buf: buffer, offset: i32) -> tuple[i32, i32]
    + returns the 1-based line and column for a byte offset
    # inspection
    -> std.text.split_lines
  playground.completions_at
    @ (buf: buffer, offset: i32, dictionary: list[string]) -> list[string]
    + returns dictionary entries that match the identifier prefix ending at offset
    + returns an empty list when the offset is not inside an identifier
    # completion
    -> std.text.prefix_match
  playground.register_executor
    @ (buf: buffer, executor: fn(string) -> result[string, string]) -> buffer
    + attaches an executor that runs snippets on demand
    # execution
  playground.run
    @ (buf: buffer) -> result[string, string]
    + returns the executor's output for the current source
    - returns error when no executor is registered
    # execution
