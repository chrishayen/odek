# Requirement: "a library for deriving a command-line interface from an arbitrary value"

Given a value (a function, a struct, or a map of commands), reflect its shape into a command tree and dispatch parsed arguments to it.

std
  std.reflect
    std.reflect.describe
      fn (value: dynamic_value) -> value_shape
      + returns the kind (function, record, map) and field or parameter info
      # reflection
    std.reflect.call
      fn (value: dynamic_value, args: list[dynamic_value]) -> result[dynamic_value, string]
      + invokes a callable value with positional arguments
      - returns error when arity does not match
      - returns error when an argument type does not match
      # reflection

auto_cli
  auto_cli.build
    fn (root: dynamic_value) -> command_tree
    + reflects a value into a command tree
    + maps record fields and map entries to subcommands
    + maps function parameters to positional arguments and flags
    # construction
    -> std.reflect.describe
  auto_cli.parse_args
    fn (tree: command_tree, argv: list[string]) -> result[parsed_command, string]
    + resolves the subcommand path and splits positional arguments from flags
    - returns error on unknown subcommand or flag
    - returns error when a required argument is missing
    # parsing
  auto_cli.dispatch
    fn (tree: command_tree, parsed: parsed_command) -> result[dynamic_value, string]
    + calls the resolved target with the parsed arguments
    - returns error when the target invocation fails
    # dispatch
    -> std.reflect.call
  auto_cli.help
    fn (tree: command_tree, path: list[string]) -> result[string, string]
    + renders a help string for the command at the given path
    - returns error when the path is unknown
    # help
  auto_cli.complete
    fn (tree: command_tree, partial: list[string]) -> list[string]
    + returns shell-style completion candidates for a partial argv
    # completion
