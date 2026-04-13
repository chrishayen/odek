# Requirement: "a source-aware minifier that understands modern language features"

A minifier pipeline that parses source to an abstract syntax tree, runs compression passes over it, and emits compact source. Only the high-level shape is fixed — individual passes are composable.

std
  std.text
    std.text.tokenize
      @ (source: string) -> result[list[token], string]
      + produces a token stream from source text
      - returns error on unterminated string or invalid character
      # lexing
  std.collections
    std.collections.map_keys
      @ (m: map[string, string]) -> list[string]
      + returns the keys of a string map in insertion order
      # collections

minifier
  minifier.parse
    @ (source: string) -> result[ast_node, string]
    + produces an AST from source text
    - returns error on syntax violations
    -> std.text.tokenize
    # parsing
  minifier.emit
    @ (node: ast_node) -> string
    + emits compact source text from an AST
    ? whitespace is kept to the minimum required for correct tokenization
    # emitting
  minifier.mangle_identifiers
    @ (node: ast_node) -> ast_node
    + renames local bindings to short identifiers while preserving scope semantics
    - leaves exported and free identifiers untouched
    # pass_mangling
  minifier.dead_code_elimination
    @ (node: ast_node) -> ast_node
    + removes unreachable statements and unused local bindings
    # pass_dce
  minifier.constant_folding
    @ (node: ast_node) -> ast_node
    + evaluates numeric and string constant expressions at compile time
    # pass_folding
  minifier.collapse_branches
    @ (node: ast_node) -> ast_node
    + eliminates branches whose condition is a compile-time constant
    # pass_branches
  minifier.inline_single_use
    @ (node: ast_node) -> ast_node
    + inlines bindings referenced exactly once when safe
    - leaves mutable or captured bindings untouched
    # pass_inlining
  minifier.rewrite_booleans
    @ (node: ast_node) -> ast_node
    + replaces true/false literals with equivalent shorter expressions
    # pass_booleans
  minifier.pipeline
    @ (passes: list[string]) -> minify_pipeline
    + returns a pipeline that runs the named passes in order
    # pipeline
  minifier.run
    @ (pipeline: minify_pipeline, source: string) -> result[string, string]
    + parses, runs all passes, then emits minified source
    - returns error when parsing fails
    # orchestration
