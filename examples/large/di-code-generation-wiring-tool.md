# Requirement: "a code generation based dependency injection wiring tool"

A library that takes provider declarations, builds a dependency graph, and emits source code that wires them into a requested type.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating or truncating it
      # filesystem

di
  di.new_graph
    fn () -> graph_state
    + creates an empty dependency graph
    # construction
  di.add_provider
    fn (state: graph_state, name: string, produces: string, depends_on: list[string]) -> result[graph_state, string]
    + registers a provider that produces a type from a list of dependency types
    - returns error when a provider with the same name already exists
    # registration
  di.set_binding
    fn (state: graph_state, interface_type: string, concrete_type: string) -> graph_state
    + maps an abstract type to the concrete type that should satisfy it
    # binding
  di.resolve_order
    fn (state: graph_state, target_type: string) -> result[list[string], string]
    + returns the provider names to invoke in order to construct the target type
    - returns error when a required type has no provider
    - returns error when the graph contains a cycle
    # resolution
  di.detect_cycles
    fn (state: graph_state) -> list[list[string]]
    + returns all dependency cycles found in the graph
    # diagnostics
  di.missing_types
    fn (state: graph_state, target_type: string) -> list[string]
    + returns types reachable from target_type that lack a provider
    # diagnostics
  di.emit_wiring_source
    fn (state: graph_state, target_type: string, function_name: string) -> result[string, string]
    + returns source code for a function that constructs the target type
    - returns error when resolve_order fails
    # code_generation
  di.parse_provider_file
    fn (source: string) -> result[list[provider_decl], string]
    + parses a provider declaration file into structured records
    - returns error when a declaration is malformed
    # parsing
  di.load_providers_from_path
    fn (state: graph_state, path: string) -> result[graph_state, string]
    + reads a provider file from disk and adds each provider to the graph
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
  di.write_generated_file
    fn (state: graph_state, target_type: string, function_name: string, path: string) -> result[void, string]
    + resolves the target and writes the emitted source to disk
    - returns error when resolution or writing fails
    # code_generation
    -> std.fs.write_all
