# Requirement: "a shell completion generator for command trees"

Given a declarative description of commands and flags, emit completion candidates for a cursor position.

std: (all units exist)

completion
  completion.new_command
    @ (name: string) -> command_node
    + creates a command node with the given name and no children
    # construction
  completion.add_subcommand
    @ (parent: command_node, child: command_node) -> command_node
    + attaches child to parent and returns the updated parent
    # tree_building
  completion.add_flag
    @ (cmd: command_node, flag: string, takes_value: bool) -> command_node
    + registers a flag on cmd, optionally expecting a value argument
    # flag_registration
  completion.resolve
    @ (root: command_node, argv: list[string]) -> command_node
    + walks argv from root and returns the deepest command reached
    - returns root when no subcommand matches the first token
    # resolution
  completion.suggest
    @ (root: command_node, argv: list[string], cursor_token: string) -> list[string]
    + returns subcommand names and flag names matching cursor_token at the resolved command
    + returns all candidates when cursor_token is empty
    # suggestion
