# Requirement: "a subcommand parser that composes with a flat flag parser"

Callers supply a flat flag parser; this library adds subcommand dispatch on top of it. Parsing returns a typed invocation; running is the caller's concern.

std: (all units exist)

subcmd
  subcmd.new
    fn (program_name: string) -> subcmd_registry
    + creates an empty registry with no commands
    # construction
  subcmd.register
    fn (r: subcmd_registry, name: string, summary: string, flag_spec: list[tuple[string, bool]]) -> subcmd_registry
    + registers a subcommand with a summary and its flag spec (name, takes_value)
    # registration
  subcmd.parse
    fn (r: subcmd_registry, argv: list[string]) -> result[subcmd_invocation, string]
    + returns the matched command, its parsed flag map, and remaining positional arguments
    - returns error when argv[0] is not a registered command
    - returns error when a flag expecting a value has none
    # parsing
  subcmd.summary
    fn (r: subcmd_registry) -> string
    + returns a one-line-per-command summary of all registered commands
    # help_text
