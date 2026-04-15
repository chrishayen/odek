# Requirement: "a JavaScript/TypeScript compiler library"

A full compilation pipeline: source text to tokens to AST to typed AST to transformed AST to output source. Large system with distinct phases.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file fully into a string
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: string) -> result[void, string]
      + writes data to path, truncating existing content
      - returns error on IO failure
      # filesystem

jsc
  jsc.tokenize
    fn (source: string) -> result[list[token], string]
    + produces tokens for identifiers, keywords, punctuation, numbers, strings, template literals, and regexes
    - returns error on unterminated strings or invalid characters, with offset
    # lexing
  jsc.parse_module
    fn (tokens: list[token]) -> result[module_ast, string]
    + builds a module AST covering statements, expressions, classes, and type annotations
    - returns error with a location on unexpected tokens
    # parsing
  jsc.parse_file
    fn (path: string) -> result[module_ast, string]
    + reads a source file and parses it
    - returns error when the file cannot be read or parsed
    # parsing
    -> std.fs.read_all
    -> jsc.tokenize
    -> jsc.parse_module
  jsc.check_types
    fn (module: module_ast) -> result[typed_module, list[diagnostic]]
    + resolves identifiers, infers types, and checks annotations
    - returns the diagnostics collected during checking on any error
    # type_checking
  jsc.strip_types
    fn (typed: typed_module) -> module_ast
    + removes all TypeScript-only syntax, leaving valid ES module AST
    # transform
  jsc.lower_esnext
    fn (module: module_ast, target: string) -> module_ast
    + lowers newer syntax (async, classes, destructuring) to the target version
    + target accepts "es5","es2015","es2017","es2020","esnext"
    # transform
  jsc.resolve_imports
    fn (module: module_ast, base_dir: string) -> result[module_ast, string]
    + rewrites relative import specifiers to canonical paths
    - returns error when a specifier cannot be resolved
    # transform
  jsc.emit_source
    fn (module: module_ast) -> string
    + prints the AST back to formatted source code
    # codegen
  jsc.emit_sourcemap
    fn (module: module_ast) -> string
    + produces a source map (v3) alongside the emitted source
    # codegen
  jsc.compile_string
    fn (source: string, target: string) -> result[string, string]
    + drives the full pipeline: tokenize, parse, check, strip, lower, emit
    - returns the first error encountered in any phase
    # pipeline
  jsc.compile_file
    fn (in_path: string, out_path: string, target: string) -> result[void, string]
    + reads a source file, compiles it, and writes the result
    - returns error on read, compile, or write failure
    # pipeline
    -> std.fs.read_all
    -> std.fs.write_all
