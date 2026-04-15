# Requirement: "a command-line flag parser with subcommand support"

Declarative command and flag registration, then parse argv into a structured result.

std: (all units exist)

flags
  flags.new_parser
    fn (name: string) -> parser
    + creates a parser with the given program name and no commands or flags
    # construction
  flags.add_command
    fn (p: parser, name: string) -> parser
    + registers a subcommand that may accept its own flags
    # command_registration
  flags.add_flag
    fn (p: parser, command: optional[string], long: string, takes_value: bool) -> parser
    + adds a flag to the root when command is absent, else to the named subcommand
    # flag_registration
  flags.parse
    fn (p: parser, argv: list[string]) -> result[parse_outcome, string]
    + returns the selected subcommand (if any), the flag map, and positional args
    - returns error on unknown flags
    - returns error when a value-taking flag is missing its value
    # parsing
  flags.get_flag
    fn (outcome: parse_outcome, name: string) -> optional[string]
    + returns the value supplied for name, or a sentinel empty string for bool flags
    # lookup
