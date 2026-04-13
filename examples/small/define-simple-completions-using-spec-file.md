# Requirement: "a library that builds shell completions from a declarative specification of commands and flags"

The project layer parses a spec and returns the set of candidate completions for a given partial input.

std
  std.encoding
    std.encoding.parse_yaml
      @ (raw: string) -> result[map[string, string], string]
      + parses a flat key/value YAML document into a string-to-string map
      - returns error on malformed YAML
      # serialization

completions
  completions.load_spec
    @ (raw: string) -> result[spec_state, string]
    + parses a spec document describing commands, subcommands, and flags
    - returns error on malformed input
    - returns error when a referenced subcommand has no definition
    # loading
    -> std.encoding.parse_yaml
  completions.suggest
    @ (spec: spec_state, tokens: list[string], prefix: string) -> list[string]
    + returns all completion candidates that match the prefix in the current command context
    + subcommand candidates are suggested before any positional argument is provided
    + flag candidates are suggested when the prefix begins with "-"
    # completion
  completions.describe
    @ (spec: spec_state, path: list[string]) -> result[string, string]
    + returns the description text for a given command path
    - returns error when the path is not defined
    # inspection
