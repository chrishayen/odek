# Requirement: "a source-to-source compiler that emits JavaScript"

A full front-end plus a single JavaScript backend. Compared to a multi-target transpiler, there is only one emitter, but the pipeline still needs lex, parse, type-check, lower, and emit.

std
  std.strings
    std.strings.split_lines
      fn (s: string) -> list[string]
      + splits s on newline characters
      # strings
    std.strings.join
      fn (parts: list[string], sep: string) -> string
      + concatenates parts with sep between
      # strings
  std.io
    std.io.read_all
      fn (path: string) -> result[string, string]
      + reads a text file into a string
      - returns error when path is unreadable
      # io
    std.io.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents to path
      - returns error when the target is not writable
      # io

js_compiler
  js_compiler.tokenize
    fn (source: string) -> result[list[token], string]
    + returns tokens for the source language
    - returns error with line and column on an unrecognized character
    # lexing
    -> std.strings.split_lines
  js_compiler.parse
    fn (tokens: list[token]) -> result[ast_program, string]
    + returns an AST program
    - returns error on unexpected tokens
    # parsing
  js_compiler.type_check
    fn (program: ast_program) -> result[typed_program, string]
    + annotates expressions with inferred types
    - returns error on type mismatches with position information
    # type_checking
  js_compiler.lower
    fn (program: typed_program) -> lowered_program
    + desugars high-level constructs (closures, channels, range loops) into simpler primitives
    # lowering
  js_compiler.emit
    fn (program: lowered_program) -> string
    + renders the lowered program as JavaScript source
    + maps source types to typed array views where possible
    # codegen
    -> std.strings.join
  js_compiler.emit_source_map
    fn (program: lowered_program) -> string
    + returns a source map JSON string linking generated lines to original positions
    # source_map
  js_compiler.compile
    fn (source: string) -> result[compile_output, string]
    + runs the full pipeline and returns generated JavaScript plus source map
    - propagates any earlier stage error
    # pipeline
  js_compiler.compile_file
    fn (source_path: string, out_path: string, map_path: string) -> result[void, string]
    + compiles source_path and writes the JavaScript and source map to disk
    - returns error on any stage failure
    # convenience
    -> std.io.read_all
    -> std.io.write_all
  js_compiler.runtime_prelude
    fn () -> string
    + returns the JavaScript runtime prelude required by emitted code (helpers for range, defer, channels)
    # runtime
