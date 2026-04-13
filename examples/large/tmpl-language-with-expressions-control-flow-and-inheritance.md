# Requirement: "a templating language with expressions, control flow, and inheritance"

Supports variable substitution, if/for blocks, filters, and template inheritance with named blocks.

std
  std.io
    std.io.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when path is not readable
      # io
  std.text
    std.text.html_escape
      @ (s: string) -> string
      + replaces &, <, >, ", ' with HTML entities
      # text
    std.text.join
      @ (parts: list[string], sep: string) -> string
      + joins parts separated by sep
      # text

tmpl
  tmpl.tokenize
    @ (source: string) -> result[list[tmpl_token], string]
    + emits text, expression ({{ ... }}), and statement ({% ... %}) tokens
    - returns error on unterminated tags
    # lexing
  tmpl.parse
    @ (tokens: list[tmpl_token]) -> result[tmpl_node, string]
    + builds a template AST with text, output, if, for, block, and extends nodes
    - returns error on unbalanced control blocks
    # parsing
  tmpl.new_env
    @ () -> env_state
    + creates an empty environment with no loaded templates and no filters
    # construction
  tmpl.register_filter
    @ (env: env_state, name: string, fn: fn(value, list[value]) -> result[value, string]) -> env_state
    + installs a named filter callable from {{ x | name(args) }}
    # filters
  tmpl.load_string
    @ (env: env_state, name: string, source: string) -> result[env_state, string]
    + parses source and stores it in the environment under name
    - returns error on parse failure
    # loading
  tmpl.load_file
    @ (env: env_state, name: string, path: string) -> result[env_state, string]
    + reads a file and loads it under name
    - returns error when the file cannot be read or parsed
    # loading
    -> std.io.read_all
  tmpl.resolve_inheritance
    @ (env: env_state, name: string) -> result[tmpl_node, string]
    + merges a template with its ancestors by overlaying child blocks onto the parent tree
    - returns error on cyclic inheritance
    # inheritance
  tmpl.render
    @ (env: env_state, name: string, context: map[string, value]) -> result[string, string]
    + executes the named template with the given context, auto-escaping output
    - returns error on undefined variables in strict mode
    # rendering
    -> std.text.html_escape
    -> std.text.join
  tmpl.render_string
    @ (env: env_state, source: string, context: map[string, value]) -> result[string, string]
    + parses and renders source in one call without naming it
    # rendering
