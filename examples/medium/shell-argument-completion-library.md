# Requirement: "a shell argument completion library"

Given a command specification and the current partial command line, produces completion candidates and renders them in the format expected by a target shell.

std: (all units exist)

completer
  completer.spec_new
    @ (name: string) -> command_spec
    + creates an empty spec for a top-level command
    # construction
  completer.add_subcommand
    @ (spec: command_spec, path: list[string], description: string) -> command_spec
    + registers a subcommand at the given path
    # spec
  completer.add_flag
    @ (spec: command_spec, path: list[string], long: string, short: string, takes_value: bool) -> command_spec
    + registers a flag scoped to a subcommand path
    # spec
  completer.add_positional
    @ (spec: command_spec, path: list[string], name: string, values: list[string]) -> command_spec
    + registers a positional argument with static value candidates
    # spec
  completer.parse_line
    @ (line: string, cursor: i32) -> tuple[list[string], string]
    + returns (tokens_before_cursor, current_token)
    ? respects single and double quotes
    # parsing
  completer.candidates
    @ (spec: command_spec, tokens: list[string], current: string) -> list[completion_candidate]
    + returns candidates matching the current partial token in context
    - returns empty list when tokens refer to an unknown subcommand
    # completion
  completer.render_bash
    @ (candidates: list[completion_candidate]) -> string
    + returns newline-separated candidates for bash-style compgen
    # rendering
  completer.render_zsh
    @ (candidates: list[completion_candidate]) -> string
    + returns candidates in zsh _describe format
    # rendering
  completer.render_fish
    @ (candidates: list[completion_candidate]) -> string
    + returns tab-separated value-description pairs
    # rendering
