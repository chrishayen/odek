# Requirement: "an interactive key-value store client shell with command autocompletion and syntax highlighting"

Reads a line from the user, tokenizes it, offers completions, highlights by token role, and dispatches to the backend. Project layer owns the shell model; std provides the terminal line reader.

std
  std.term
    std.term.read_line
      fn (prompt: string) -> result[string, string]
      + reads one line of input, showing the prompt
      - returns error on end-of-input
      # terminal
    std.term.write
      fn (text: string) -> void
      + writes text to the terminal
      # terminal

kv_shell
  kv_shell.new
    fn (commands: list[command_spec]) -> shell_state
    + creates a shell with the given command registry
    + each command_spec has name, arity, and description
    # construction
  kv_shell.tokenize
    fn (line: string) -> list[token]
    + splits the input line into tokens preserving quoted strings
    # parsing
  kv_shell.complete
    fn (shell: shell_state, prefix: string) -> list[string]
    + returns command names and argument hints matching prefix
    + returns all commands when prefix is empty
    # completion
  kv_shell.highlight
    fn (shell: shell_state, line: string) -> list[highlight_span]
    + annotates tokens as command, string, number, or error
    # rendering
  kv_shell.render_highlighted
    fn (spans: list[highlight_span]) -> string
    + returns the line with ansi color escapes applied to each span
    # rendering
  kv_shell.execute
    fn (shell: shell_state, line: string, backend: backend_handle) -> result[string, string]
    + parses the line, validates arity, runs the command, and returns its reply text
    - returns error when the command is unknown
    - returns error when argument count does not match the spec
    # dispatch
  kv_shell.run_repl
    fn (shell: shell_state, backend: backend_handle) -> result[void, string]
    + runs the read-eval-print loop until EOF
    # control
    -> std.term.read_line
    -> std.term.write
