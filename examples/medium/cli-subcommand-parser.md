# Requirement: "a subcommand-oriented command-line argument parser"

A parser in the spirit of tool-with-subcommands style interfaces. Callers register commands with flags; the parser dispatches based on argv. No execution wiring — the parser returns a structured result and the caller runs the command.

std: (all units exist)

cli_parser
  cli_parser.new
    fn (program_name: string) -> parser_state
    + creates an empty parser with no registered commands
    # construction
  cli_parser.register_command
    fn (state: parser_state, name: string, summary: string) -> parser_state
    + registers a top-level subcommand with a short summary
    # registration
  cli_parser.register_flag
    fn (state: parser_state, command: string, flag: string, takes_value: bool) -> result[parser_state, string]
    + attaches a flag to the named command
    - returns error when the command has not been registered
    # registration
  cli_parser.parse
    fn (state: parser_state, argv: list[string]) -> result[parsed_invocation, string]
    + returns the matched command name, its flag map, and positional arguments
    - returns error when argv names an unknown command
    - returns error when a flag that requires a value is missing one
    # parsing
  cli_parser.render_help
    fn (state: parser_state) -> string
    + returns a human-readable help summary listing commands and their flags
    # help_text
