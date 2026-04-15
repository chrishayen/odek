# Requirement: "a hierarchical command-line interaction framework"

Users register named subcommands with flags, then dispatch an argv vector to the right handler. The framework owns parsing and lookup; handlers are opaque callables.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s on sep and returns all segments
      # strings
    std.strings.has_prefix
      fn (s: string, prefix: string) -> bool
      + returns true when s starts with prefix
      # strings

command_tree
  command_tree.new
    fn (program_name: string) -> command_state
    + creates an empty root command with the given program name
    # construction
  command_tree.register
    fn (state: command_state, path: list[string], handler_id: string, flags: list[string]) -> result[command_state, string]
    + attaches a subcommand at path (e.g. ["server","start"]) bound to handler_id
    + declares the flags that subcommand accepts
    - returns error when path already exists
    # registration
  command_tree.parse
    fn (state: command_state, argv: list[string]) -> result[parsed_invocation, string]
    + walks argv, matches the deepest registered subcommand, and returns handler_id plus parsed flags and positional args
    - returns error when the first non-flag token is not a registered subcommand
    - returns error when an unknown flag is supplied
    # parsing
    -> std.strings.has_prefix
  command_tree.usage
    fn (state: command_state, path: list[string]) -> result[string, string]
    + returns a human-readable help string for the subcommand at path
    - returns error when path is not registered
    # help
    -> std.strings.split
  command_tree.dispatch
    fn (state: command_state, argv: list[string], handlers: map[string, command_handler]) -> result[i32, string]
    + parses argv and invokes the bound handler, returning its exit code
    - returns error when no handler is registered for the resolved handler_id
    # dispatch
