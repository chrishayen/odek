# Requirement: "a static analyzer that detects nested context values inside loops or closures"

Walks a parsed syntax tree, tracks context variable assignments, and flags ones that shadow a parent context inside a loop or closure body.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the entire file as a string
      - returns error when the file does not exist
      # io
  std.ast
    std.ast.parse_source
      @ (source: string) -> result[ast_node, string]
      + parses source text into a generic syntax tree
      - returns error on syntax error
      # parsing
    std.ast.walk
      @ (root: ast_node, visit: fn(ast_node) -> bool) -> void
      + invokes visit for every descendant node
      ? visit returning false skips the subtree
      # traversal

fatcontext
  fatcontext.find_context_assignments
    @ (tree: ast_node) -> list[assignment]
    + returns every assignment whose right side wraps another context variable
    # scanning
    -> std.ast.walk
  fatcontext.is_inside_loop_or_closure
    @ (tree: ast_node, target: ast_node) -> bool
    + returns true when target has a loop or closure ancestor
    - returns false for top-level assignments
    # scope_analysis
  fatcontext.diagnose
    @ (tree: ast_node) -> list[diagnostic]
    + returns diagnostics for each nested context assignment inside a loop or closure
    + each diagnostic carries line and column from the source
    # reporting
  fatcontext.analyze_file
    @ (path: string) -> result[list[diagnostic], string]
    + parses the file and returns all diagnostics
    - returns error when parsing fails
    # entry
    -> std.fs.read_all
    -> std.ast.parse_source
