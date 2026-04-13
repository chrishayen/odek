# Requirement: "a documentation generator for environment variable declarations"

Given parsed source that describes the environment variables a program reads, render a reference document.

std: (all units exist)

envdoc
  envdoc.new_variable
    @ (name: string, type_name: string, default_value: string, description: string, required: bool) -> env_var
    + builds an env var record
    # model
  envdoc.parse_comment_block
    @ (raw: string) -> result[env_var, string]
    + parses a doc comment of the form "ENV_NAME (type, default=...): description" into an env_var
    - returns error when the header line is missing or malformed
    # parsing
  envdoc.render_markdown
    @ (vars: list[env_var]) -> string
    + returns a markdown table listing each variable, its type, default, required flag, and description
    ? variables are rendered in the order supplied
    # rendering
