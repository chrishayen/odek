# Requirement: "a command-line framework with argument parsing, documentation generation, and a plugin system"

Commands are declared with typed flags, parsed into a typed bag, and can be documented automatically from their declarations. A plugin registry lets external modules add commands.

std
  std.strings
    std.strings.starts_with
      fn (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
    std.strings.join
      fn (parts: list[string], sep: string) -> string
      + joins parts with sep between each
      # strings

cli
  cli.command_new
    fn (name: string, summary: string) -> command_spec
    + constructs a command with an empty flag set
    # definition
  cli.command_flag
    fn (spec: command_spec, name: string, kind: flag_kind, required: bool, help: string) -> command_spec
    + appends a flag declaration to the command spec
    # definition
  cli.command_arg
    fn (spec: command_spec, name: string, kind: flag_kind, help: string) -> command_spec
    + appends a positional argument declaration
    # definition
  cli.registry_new
    fn () -> registry_state
    + constructs an empty command registry
    # registration
  cli.registry_add
    fn (state: registry_state, spec: command_spec, handler: command_fn) -> registry_state
    + registers a command and its handler
    # registration
  cli.registry_add_plugin
    fn (state: registry_state, plugin: plugin) -> registry_state
    + merges every command contributed by a plugin into the registry
    # plugins
  cli.parse
    fn (state: registry_state, args: list[string]) -> result[parsed_invocation, string]
    + resolves the command and parses flags and positionals per its spec
    - returns error when the command is unknown
    - returns error when a required flag is absent
    - returns error when a flag value cannot be coerced to its declared kind
    # parsing
    -> std.strings.starts_with
  cli.render_help
    fn (spec: command_spec) -> string
    + renders a human-readable help page from a command spec
    # documentation
    -> std.strings.join
  cli.render_index
    fn (state: registry_state) -> string
    + renders a listing of every registered command with its summary
    # documentation
    -> std.strings.join
  cli.dispatch
    fn (state: registry_state, args: list[string]) -> result[i32, string]
    + parses args and invokes the matched handler, returning its exit code
    - returns error when parsing fails
    # execution
    -> cli.parse
