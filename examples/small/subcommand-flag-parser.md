# Requirement: "a command-line flag parser that supports subcommands"

Defines a root flag set with subcommands, each with its own flags, and parses an argument vector against them.

std: (all units exist)

subflag
  subflag.new
    fn (name: string) -> parser_state
    + creates a parser with the given program name and no flags or subcommands
    # construction
  subflag.add_flag
    fn (parser: parser_state, sub: string, name: string, default_value: string) -> parser_state
    + registers a flag on the given subcommand, or on the root when sub is ""
    # registration
  subflag.add_subcommand
    fn (parser: parser_state, name: string) -> parser_state
    + registers a subcommand name
    # registration
  subflag.parse
    fn (parser: parser_state, args: list[string]) -> result[parse_result, string]
    + returns the selected subcommand, flag values, and positional arguments
    - returns error when the first non-flag argument is not a registered subcommand
    - returns error on unknown flags
    # parsing
  subflag.usage
    fn (parser: parser_state) -> string
    + returns a human-readable usage summary listing subcommands and their flags
    # help
