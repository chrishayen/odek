# Requirement: "a source-to-source compiler with multiple target backends"

Parses a source language into an AST, lowers it to a backend-agnostic IR, then emits textual output for one of several target languages. Each backend is a pure IR-to-string function.

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

transpiler
  transpiler.tokenize
    fn (source: string) -> result[list[token], string]
    + returns tokens for the source language
    - returns error with line and column on an unrecognized character
    # lexing
    -> std.strings.split_lines
  transpiler.parse
    fn (tokens: list[token]) -> result[ast_program, string]
    + returns an AST program: declarations, functions, statements
    - returns error on unexpected tokens
    # parsing
  transpiler.lower_to_ir
    fn (program: ast_program) -> result[ir_program, string]
    + converts the AST into a target-independent IR with typed SSA-style instructions
    - returns error when a language construct has no IR equivalent
    # lowering
  transpiler.optimize
    fn (program: ir_program) -> ir_program
    + applies constant folding and dead-store elimination
    # optimization
  transpiler.emit_javascript
    fn (program: ir_program) -> string
    + renders program as JavaScript source
    # backend_js
    -> std.strings.join
  transpiler.emit_cpp
    fn (program: ir_program) -> string
    + renders program as C++ source
    # backend_cpp
    -> std.strings.join
  transpiler.emit_java
    fn (program: ir_program) -> string
    + renders program as Java source
    # backend_java
    -> std.strings.join
  transpiler.emit_csharp
    fn (program: ir_program) -> string
    + renders program as C# source
    # backend_cs
    -> std.strings.join
  transpiler.compile
    fn (source: string, target: string) -> result[string, string]
    + full pipeline: tokenize, parse, lower, optimize, emit for target in {"javascript","cpp","java","csharp"}
    - returns error when target is not one of the supported names
    - propagates any earlier stage error
    # pipeline
  transpiler.compile_file
    fn (source_path: string, target: string, out_path: string) -> result[void, string]
    + reads source_path, compiles to target, and writes the result to out_path
    - returns error on any stage failure
    # convenience
    -> std.io.read_all
    -> std.io.write_all
