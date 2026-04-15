# Requirement: "a command-line argument parser"

Declarative: the caller builds a spec of flags and positionals, then parses an argv list against it.

std: (all units exist)

cli_parser
  cli_parser.new_spec
    fn (program_name: string) -> cli_spec
    + creates an empty argument spec with the given program name
    # construction
  cli_parser.add_flag
    fn (spec: cli_spec, long_name: string, short_name: string, help: string) -> cli_spec
    + registers a boolean flag
    # construction
  cli_parser.add_option
    fn (spec: cli_spec, long_name: string, short_name: string, default_value: string, help: string) -> cli_spec
    + registers a string-valued option with a default
    # construction
  cli_parser.add_positional
    fn (spec: cli_spec, name: string, required: bool, help: string) -> cli_spec
    + appends a positional argument to the spec in declaration order
    # construction
  cli_parser.parse
    fn (spec: cli_spec, argv: list[string]) -> result[parsed_args, string]
    + returns parsed flags, options, and positionals when argv matches the spec
    + accepts both "--name value" and "--name=value" forms
    + treats "--" as the end of option parsing
    - returns error when a required positional is missing
    - returns error when an unknown option is supplied
    # parsing
  cli_parser.get_flag
    fn (args: parsed_args, long_name: string) -> bool
    + returns whether the named flag was present
    # access
  cli_parser.get_option
    fn (args: parsed_args, long_name: string) -> string
    + returns the option value, falling back to its declared default
    # access
  cli_parser.get_positional
    fn (args: parsed_args, name: string) -> optional[string]
    + returns the positional value by declared name or none
    # access
  cli_parser.render_help
    fn (spec: cli_spec) -> string
    + returns a formatted usage string listing the program, options, and positionals
    # help
