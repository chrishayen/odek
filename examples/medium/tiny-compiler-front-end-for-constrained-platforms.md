# Requirement: "a compiler front-end that targets constrained platforms"

A minimal pipeline that parses source text, produces an intermediate representation, and emits for a chosen backend.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the contents of the file at path
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to the file at path, creating it if missing
      # filesystem

tiny_compiler
  tiny_compiler.tokenize
    fn (source: string) -> result[list[token], string]
    + returns a token stream from source text
    - returns error on an unterminated string literal
    # lexing
  tiny_compiler.parse
    fn (tokens: list[token]) -> result[ast_node, string]
    + returns an AST root node from a token stream
    - returns error on unexpected token
    # parsing
  tiny_compiler.lower
    fn (ast: ast_node) -> result[ir_module, string]
    + lowers the AST into an intermediate representation module
    - returns error on unsupported language construct
    # lowering
  tiny_compiler.emit
    fn (ir: ir_module, target: string) -> result[bytes, string]
    + returns the encoded artifact for the named target backend
    - returns error when the target is not supported
    # code_generation
  tiny_compiler.compile_file
    fn (source_path: string, output_path: string, target: string) -> result[void, string]
    + reads source from disk, runs the pipeline, and writes the artifact
    - returns error when any stage fails
    # pipeline
    -> std.fs.read_all
    -> std.fs.write_all
