# Requirement: "a library that starts an interactive REPL with access to a captured variable scope"

The caller provides a snapshot of named values; the REPL lets users read and mutate them through line-based input.

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

scope_repl
  scope_repl.session_new
    @ (scope: map[string, value]) -> repl_state
    + creates a session backed by the given variable scope
    # construction
  scope_repl.evaluate_line
    @ (s: repl_state, line: string) -> result[string, string]
    + interprets "name" as read and "name = literal" as assignment
    + returns the rendered value or confirmation message
    - returns error when name is not in the scope
    - returns error on an unparseable literal on the right-hand side
    # evaluation
  scope_repl.run
    @ (s: repl_state, prompt: string) -> result[void, string]
    + reads lines until eof, printing the result of each evaluate_line
    + stops when the user enters a bare "exit" line
    # loop
    -> std.io.read_line
    -> std.io.write_string
