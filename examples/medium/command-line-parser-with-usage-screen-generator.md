# Requirement: "a command-line parser that generates a usage screen from its declaration"

Parse argv into a structured result and format a usage message from the same declaration.

std: (all units exist)

argparse
  argparse.new
    fn (program: string, description: string) -> parser
    + creates an empty parser with metadata for usage rendering
    # construction
  argparse.option
    fn (p: parser, long: string, help: string, takes_value: bool) -> parser
    + registers a named option
    # option_registration
  argparse.positional
    fn (p: parser, name: string, help: string) -> parser
    + registers a required positional argument in order
    # positional_registration
  argparse.parse
    fn (p: parser, argv: list[string]) -> result[parse_result, string]
    + returns values for every registered option and positional
    - returns error when a required positional is missing
    - returns error on unknown options
    # parsing
  argparse.usage
    fn (p: parser) -> string
    + renders a multi-line usage screen listing options and positionals with help text
    # usage_rendering
