# Requirement: "a code generator that produces mock method implementations for embedded composite types"

Walks type declarations, finds methods on an embedded parent type, and emits overridable mock stubs that record calls and return configured values.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns file contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, content: string) -> result[void, string]
      + creates or overwrites the target file
      # filesystem

moxie
  moxie.parse_source
    fn (source: string) -> result[source_ast, string]
    + parses type declarations, methods, and embedded composition edges
    - returns error on syntax problems
    # parsing
  moxie.find_embedded_methods
    fn (ast: source_ast, host_type: string) -> list[method_sig]
    + returns methods inherited via embedding from any parent type
    ? resolves only one level of embedding
    # inheritance_walk
  moxie.render_call_recorder
    fn (type_name: string, methods: list[method_sig]) -> string
    + emits a recorder struct with one slice-of-calls field per method
    # codegen
  moxie.render_mock_methods
    fn (type_name: string, methods: list[method_sig]) -> string
    + emits each method as a stub that records its arguments and returns a configurable value
    # codegen
  moxie.generate
    fn (src_path: string, host_type: string, out_path: string) -> result[void, string]
    + reads source, finds embedded methods, renders mock, writes it
    - returns error when host_type is not found in the source
    # orchestration
    -> std.fs.read_all
    -> std.fs.write_all
