# Requirement: "a composable command-line interface definition library"

Declares commands, options, and arguments, then parses a raw argv list into a typed invocation. No I/O: printing and exiting are the caller's job.

std
  std.strings
    std.strings.split
      fn (input: string, sep: string) -> list[string]
      + splits on every occurrence of sep
      + returns a single-element list when sep is absent
      # strings
    std.strings.has_prefix
      fn (input: string, prefix: string) -> bool
      + reports whether input begins with prefix
      # strings

cli
  cli.new_command
    fn (name: string, summary: string) -> command_def
    + creates a command with the given name and short description
    # definition
  cli.add_option
    fn (cmd: command_def, long: string, short: string, kind: string, default_value: string) -> command_def
    + registers a named option with a type tag and default
    ? kind is one of "string", "int", "bool"
    # definition
  cli.add_argument
    fn (cmd: command_def, name: string, required: bool) -> command_def
    + adds a positional argument after existing ones
    # definition
  cli.add_subcommand
    fn (parent: command_def, child: command_def) -> command_def
    + nests child under parent so parsing dispatches by name
    # composition
  cli.parse
    fn (cmd: command_def, argv: list[string]) -> result[invocation, string]
    + returns an invocation with resolved options, positionals, and subcommand path
    - returns error on unknown options
    - returns error when a required argument is missing
    # parsing
    -> std.strings.has_prefix
    -> std.strings.split
  cli.help_text
    fn (cmd: command_def) -> string
    + renders usage, options, and subcommands as formatted text
    # help
