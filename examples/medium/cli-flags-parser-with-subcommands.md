# Requirement: "a command-line flags parser with subcommand support"

Declarative parser: the caller builds a command tree, then feeds in argv and receives a parsed result.

std: (all units exist)

flags
  flags.new_command
    @ (name: string) -> command_state
    + creates a command with the given name, no flags, and no subcommands
    # construction
  flags.add_string
    @ (cmd: command_state, name: string, default_value: string) -> command_state
    + registers a string-valued flag
    - returns unchanged command when name is empty
    # flag_registration
  flags.add_int
    @ (cmd: command_state, name: string, default_value: i64) -> command_state
    + registers an integer-valued flag
    # flag_registration
  flags.add_bool
    @ (cmd: command_state, name: string) -> command_state
    + registers a boolean switch that defaults to false
    # flag_registration
  flags.add_subcommand
    @ (cmd: command_state, child: command_state) -> command_state
    + attaches a subcommand under the parent
    # subcommand
  flags.parse
    @ (cmd: command_state, argv: list[string]) -> result[parsed_args, string]
    + returns the matched command path, flag values, and positional args
    - returns error when a required flag value is missing
    - returns error when an unknown flag appears and strict mode is set
    # parsing
  flags.usage
    @ (cmd: command_state) -> string
    + returns a usage string listing flags, subcommands, and positional arguments
    # help
