# Requirement: "an interactive SQL shell library with autocompletion and syntax highlighting"

A library that drives a line-based interactive SQL session. Input/output are abstracted so the caller supplies the terminal.

std
  std.io
    std.io.read_line
      @ () -> result[string, string]
      + returns the next line from stdin without the trailing newline
      - returns error on eof
      # io
    std.io.write_string
      @ (s: string) -> void
      + writes s to stdout
      # io

sqlshell
  sqlshell.session_new
    @ (schema: list[table_info]) -> session_state
    + creates a session that knows the given tables and their columns
    # construction
  sqlshell.tokenize
    @ (input: string) -> list[sql_token]
    + splits input into keyword, identifier, literal, punctuation tokens with byte offsets
    # lexing
  sqlshell.highlight
    @ (tokens: list[sql_token]) -> string
    + returns input with ANSI color codes wrapping each token by category
    # rendering
  sqlshell.complete
    @ (s: session_state, input: string, cursor: i32) -> list[string]
    + returns candidate completions at the cursor position
    + suggests table names after FROM and column names after SELECT or WHERE
    - returns [] when the cursor is inside a string literal
    # completion
  sqlshell.execute_line
    @ (s: session_state, line: string, run: fn(string) -> result[list[row], string]) -> result[list[row], string]
    + sends the line to run when it ends with a semicolon
    + accumulates lines in session buffer until a complete statement is formed
    - returns error from run unchanged
    # execution
  sqlshell.format_rows
    @ (rows: list[row]) -> string
    + renders rows as a box-drawing table with column headers
    + returns "(no rows)" for an empty result
    # formatting
