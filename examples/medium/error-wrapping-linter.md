# Requirement: "a linter that checks returned errors from external modules are wrapped"

Parses source, walks call sites, and reports unwrapped external errors as diagnostics. Language parsing is a std primitive so the linter focuses on the policy.

std
  std.source
    std.source.parse
      @ (src: string) -> result[ast_node, string]
      + parses source text into an abstract syntax tree
      - returns error on syntax failure
      # parsing
    std.source.walk
      @ (root: ast_node) -> list[ast_node]
      + returns every node in pre-order
      # traversal

wrapcheck
  wrapcheck.analyze
    @ (src: string, local_modules: list[string]) -> result[list[diagnostic], string]
    + returns a diagnostic for each call that returns an error from a non-local module without wrapping
    - returns error when source fails to parse
    # analysis
    -> std.source.parse
    -> std.source.walk
  wrapcheck.is_external_call
    @ (node: ast_node, local_modules: list[string]) -> bool
    + returns true when the call target's module is not in local_modules
    # classification
  wrapcheck.is_wrapped_return
    @ (node: ast_node) -> bool
    + returns true when the error flows through a wrapping call before being returned
    # policy
  wrapcheck.format_diagnostic
    @ (diag: diagnostic) -> string
    + formats a diagnostic as "file:line:col: message"
    # reporting
