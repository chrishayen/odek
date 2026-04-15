# Requirement: "a command-line option and argument parser"

Declarative parser definition plus a parse step that turns argv into typed values and generates help text.

std: (all units exist)

argparse
  argparse.new_parser
    fn (program: string, description: string) -> parser_state
    + creates a parser with the given program name and description
    # construction
  argparse.add_flag
    fn (state: parser_state, long: string, short: string, help: string) -> parser_state
    + adds a boolean flag toggled by --long or -short
    # registration
  argparse.add_option
    fn (state: parser_state, long: string, short: string, kind: value_kind, default: optional[arg_value], help: string) -> parser_state
    + adds a named option that consumes one value of the given kind
    # registration
  argparse.add_positional
    fn (state: parser_state, name: string, kind: value_kind, required: bool, help: string) -> parser_state
    + adds a positional argument
    - panics when a required positional follows an optional one
    # registration
  argparse.add_subcommand
    fn (state: parser_state, name: string, sub: parser_state) -> parser_state
    + attaches a nested parser under the given subcommand name
    # registration
  argparse.parse
    fn (state: parser_state, argv: list[string]) -> result[parsed_args, string]
    + returns a map of argument names to values
    - returns error on unknown option
    - returns error on missing required positional
    - returns error on type conversion failure
    # parsing
  argparse.format_help
    fn (state: parser_state) -> string
    + returns a formatted usage and options block
    # help
