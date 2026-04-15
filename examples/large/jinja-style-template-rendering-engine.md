# Requirement: "a template rendering engine with jinja-style syntax"

Parses a template into tokens and an AST, then renders against a context. Supports variable interpolation, conditionals, loops, and filters.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when path does not exist
      # filesystem
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s by sep
      # strings
    std.strings.trim
      fn (s: string) -> string
      + removes leading and trailing whitespace
      # strings

template_engine
  template_engine.tokenize
    fn (source: string) -> result[list[token], string]
    + splits source into text, expression ({{ }}), and statement ({% %}) tokens
    - returns error on unclosed delimiter
    # lexing
  template_engine.parse
    fn (tokens: list[token]) -> result[ast_node, string]
    + builds a tree of text, interpolation, if, and for nodes
    - returns error when if has no matching endif
    - returns error when for has no matching endfor
    # parsing
  template_engine.eval_expression
    fn (expr: string, context: map[string,string]) -> result[string, string]
    + evaluates a dotted-path lookup against context
    + supports pipe filters: "name|upper"
    - returns error on unknown filter
    - returns error on unresolved path
    # expression_evaluation
  template_engine.register_filter
    fn (engine: engine_state, name: string, fn: filter_fn) -> engine_state
    + adds a named filter usable in expressions
    # configuration
  template_engine.new_engine
    fn () -> engine_state
    + creates an engine with the default filter set (upper, lower, length)
    # construction
  template_engine.render
    fn (engine: engine_state, source: string, context: map[string,string]) -> result[string, string]
    + tokenizes, parses, and evaluates source against context
    - returns error when any stage fails
    # rendering
  template_engine.render_file
    fn (engine: engine_state, path: string, context: map[string,string]) -> result[string, string]
    + reads a template file and renders it
    - returns error when the file cannot be read
    # rendering
    -> std.fs.read_all
  template_engine.extend_template
    fn (child: string, parent: string) -> result[string, string]
    + resolves {% block %} overrides from child into parent
    - returns error when a block name does not exist in parent
    # inheritance
  template_engine.execute_for
    fn (var_name: string, iterable: list[string], body: ast_node, context: map[string,string]) -> result[string, string]
    + iterates the body with var_name bound to each element
    # control_flow
  template_engine.execute_if
    fn (condition: string, then_branch: ast_node, else_branch: optional[ast_node], context: map[string,string]) -> result[string, string]
    + renders the branch whose condition evaluates truthy
    # control_flow
