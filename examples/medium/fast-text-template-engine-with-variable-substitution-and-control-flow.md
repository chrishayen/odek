# Requirement: "a fast text template engine with variable substitution and control flow"

Compiles templates into an intermediate form and renders them against a context map.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads an entire file as text
      - returns error when the path cannot be opened
      # filesystem

template_engine
  template_engine.parse
    fn (source: string) -> result[template, string]
    + parses a template source into an AST of text, expressions, and blocks
    - returns error on unbalanced braces or unknown directives
    # parsing
  template_engine.parse_file
    fn (path: string) -> result[template, string]
    + reads and parses a template from disk
    - returns error when the file is missing or malformed
    # parsing
    -> std.fs.read_all
    -> template_engine.parse
  template_engine.compile
    fn (tpl: template) -> compiled_template
    + lowers a parsed template into a flat instruction stream for fast rendering
    # compilation
  template_engine.render
    fn (tpl: compiled_template, context: map[string, string]) -> result[string, string]
    + renders a compiled template by substituting variables from the context
    + supports if and for blocks against context entries
    - returns error when a referenced variable is missing
    # rendering
  template_engine.render_source
    fn (source: string, context: map[string, string]) -> result[string, string]
    + parses, compiles, and renders a template in one call
    - returns error at any stage
    # rendering
    -> template_engine.parse
    -> template_engine.compile
    -> template_engine.render
