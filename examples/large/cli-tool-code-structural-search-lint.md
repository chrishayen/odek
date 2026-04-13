# Requirement: "a code structural search, lint and rewrite library"

Match source code against AST patterns, report findings, and rewrite matched nodes. Parsing is delegated to a pluggable grammar.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full file contents as a string
      - returns error when the path is unreadable
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + writes content to the path atomically
      - returns error on write failure
      # filesystem
  std.fs.walk
    std.fs.walk.list_files
      @ (root: string, extensions: list[string]) -> result[list[string], string]
      + returns all file paths under root matching the extensions
      - returns error when root does not exist
      # filesystem

code_search
  code_search.parse_source
    @ (source: string, language: string) -> result[ast_node, string]
    + returns the root AST node for valid source
    - returns error on syntax errors with line and column
    # parsing
  code_search.parse_pattern
    @ (pattern: string, language: string) -> result[ast_pattern, string]
    + returns a pattern that can contain metavariables like $X
    - returns error when the pattern is not a valid fragment
    # pattern
  code_search.match_node
    @ (node: ast_node, pattern: ast_pattern) -> optional[map[string, ast_node]]
    + returns metavariable bindings when the node matches
    - returns none when the node does not match
    # matching
  code_search.find_all
    @ (root: ast_node, pattern: ast_pattern) -> list[ast_match]
    + returns every match found anywhere in the subtree
    + empty list when nothing matches
    # search
  code_search.rewrite
    @ (root: ast_node, pattern: ast_pattern, replacement: string) -> result[string, string]
    + returns new source with all matches replaced using metavariable bindings
    - returns error when the replacement references an unbound metavariable
    # rewrite
  code_search.load_rule
    @ (raw: string) -> result[lint_rule, string]
    + parses a rule with id, message, pattern, and severity fields
    - returns error when required fields are missing
    # rules
  code_search.apply_rule
    @ (rule: lint_rule, source: string) -> result[list[lint_finding], string]
    + returns findings with line, column, and the rule message
    - returns error when the source fails to parse
    # linting
  code_search.scan_tree
    @ (root_path: string, rules: list[lint_rule]) -> result[list[lint_finding], string]
    + walks the tree and applies every rule to every file
    + aggregates findings across files
    # scanning
    -> std.fs.walk.list_files
    -> std.fs.read_all
  code_search.apply_fix
    @ (source: string, pattern: ast_pattern, replacement: string) -> result[string, string]
    + returns rewritten source
    - returns error when the source fails to parse
    # rewrite
  code_search.write_fixed_file
    @ (path: string, new_source: string) -> result[void, string]
    + writes the rewritten source back to disk
    - returns error on write failure
    # filesystem
    -> std.fs.write_all
