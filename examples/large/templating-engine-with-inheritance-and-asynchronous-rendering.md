# Requirement: "a templating engine with inheritance and asynchronous rendering"

Templates can extend a parent and override named blocks. Rendering returns a future so expressions may resolve asynchronously.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when path does not exist
      # filesystem
  std.async
    std.async.resolve
      @ (value: string) -> future[string]
      + creates an immediately-ready future
      # async
    std.async.all
      @ (futures: list[future[string]]) -> future[list[string]]
      + joins a list of futures into one
      # async

templating
  templating.parse
    @ (source: string) -> result[template_ast, string]
    + parses source into nodes: text, variable, block, extends
    - returns error on unterminated tag
    # parsing
  templating.resolve_inheritance
    @ (child: template_ast, parent: template_ast) -> result[template_ast, string]
    + merges child block overrides into parent
    - returns error when child overrides a block the parent does not declare
    # inheritance
  templating.render_async
    @ (ast: template_ast, context: map[string,string]) -> future[result[string, string]]
    + evaluates nodes concurrently and returns a future of the rendered string
    - future resolves to error when any variable is unresolved
    # rendering
    -> std.async.resolve
    -> std.async.all
  templating.register_async_helper
    @ (engine: engine_state, name: string, handler: async_helper_fn) -> engine_state
    + adds a helper whose value resolves via a future
    # configuration
  templating.new_engine
    @ () -> engine_state
    + creates an engine with no registered helpers
    # construction
  templating.load_template
    @ (engine: engine_state, path: string) -> result[template_ast, string]
    + reads and parses a template from disk
    - returns error when the file is missing
    # loading
    -> std.fs.read_all
  templating.render_file_async
    @ (engine: engine_state, path: string, context: map[string,string]) -> future[result[string, string]]
    + loads, resolves inheritance, and renders in one call
    # rendering
  templating.define_block
    @ (ast: template_ast, name: string, body: template_ast) -> template_ast
    + attaches a named block override to a child template
    # inheritance
