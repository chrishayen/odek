# Requirement: "a pluggable source-code linting utility"

A linter that parses source into an AST, walks it with pluggable rules, and reports diagnostics.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the full file as UTF-8
      - returns error when the file does not exist
      # filesystem

linter
  linter.parse_source
    @ (source: string) -> result[ast_node, string]
    + returns the root AST node for the source
    - returns error with line and column on parse failure
    # parsing
  linter.new_registry
    @ () -> registry_state
    + returns an empty rule registry
    # construction
  linter.register_rule
    @ (state: registry_state, rule: rule_def) -> registry_state
    + adds a rule definition keyed by its id
    # registration
  linter.load_config
    @ (raw: string) -> result[config_state, string]
    + parses a configuration document that enables rules and sets severities
    - returns error when an unknown severity is set
    # configuration
  linter.lint_source
    @ (registry: registry_state, config: config_state, source: string) -> result[list[diagnostic], string]
    + returns diagnostics for every enabled rule that fires
    - returns error when source cannot be parsed
    # analysis
  linter.lint_file
    @ (registry: registry_state, config: config_state, path: string) -> result[list[diagnostic], string]
    + reads path and lints its contents
    - returns error when the file cannot be read
    # analysis
    -> std.fs.read_all
  linter.walk_ast
    @ (root: ast_node, visitor: visitor_fn) -> void
    + invokes visitor for every node in depth-first order
    # traversal
  linter.report_text
    @ (diags: list[diagnostic]) -> string
    + returns a human-readable report grouped by file
    # reporting
  linter.report_json
    @ (diags: list[diagnostic]) -> string
    + returns diagnostics encoded as a JSON array
    # reporting
  linter.has_errors
    @ (diags: list[diagnostic]) -> bool
    + returns true when any diagnostic has severity error
    # classification
  linter.autofix
    @ (source: string, diags: list[diagnostic]) -> result[string, string]
    + returns source with non-overlapping fixes applied
    - returns error when fixes conflict
    # autofix
