# Requirement: "a declarative command-line parser driven by structured field descriptors"

A caller describes commands and flags as data; the library binds argv to typed values. No code generation, no reflection on a host language — the schema is explicit.

std
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits on every occurrence of the separator
      # strings
    std.strings.has_prefix
      @ (s: string, prefix: string) -> bool
      + returns true when s starts with prefix
      # strings
  std.parse
    std.parse.to_i64
      @ (s: string) -> result[i64, string]
      + parses a decimal integer, allowing an optional leading sign
      - returns error on non-numeric input
      # parsing
    std.parse.to_f64
      @ (s: string) -> result[f64, string]
      + parses a decimal float
      - returns error on non-numeric input
      # parsing

cli_schema
  cli_schema.new_command
    @ (name: string, summary: string) -> command_spec
    + creates an empty command with no flags or subcommands
    # construction
  cli_schema.add_flag
    @ (cmd: command_spec, long: string, short: string, kind: string, required: bool, default_value: string) -> command_spec
    + registers a flag under the given command
    + kind is one of "string", "i64", "f64", "bool"
    - returns an error-marked command_spec when long starts with a dash
    # schema
  cli_schema.add_subcommand
    @ (parent: command_spec, child: command_spec) -> command_spec
    + attaches a child command under the parent
    # schema
  cli_schema.parse
    @ (root: command_spec, argv: list[string]) -> result[parsed_invocation, string]
    + resolves the subcommand chain and collects typed flag values and positional args
    + supports "--flag value", "--flag=value", and short "-f value" forms
    - returns error when a required flag is missing
    - returns error when an unknown flag is supplied
    - returns error when a flag value fails to parse for its declared kind
    # parsing
    -> std.strings.split
    -> std.strings.has_prefix
    -> std.parse.to_i64
    -> std.parse.to_f64
  cli_schema.render_help
    @ (cmd: command_spec) -> string
    + produces a help text listing usage, flags, and subcommands
    # help
  cli_schema.get_string
    @ (invocation: parsed_invocation, name: string) -> optional[string]
    + retrieves a string flag value by its long name
    # accessor
  cli_schema.get_i64
    @ (invocation: parsed_invocation, name: string) -> optional[i64]
    + retrieves an integer flag value by its long name
    # accessor
  cli_schema.get_bool
    @ (invocation: parsed_invocation, name: string) -> bool
    + returns true when the named boolean flag was set
    # accessor
