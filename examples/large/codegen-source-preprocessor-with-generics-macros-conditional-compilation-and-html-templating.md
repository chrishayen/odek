# Requirement: "a source preprocessor with generics, macros, conditional compilation, and HTML templating"

A library that takes annotated source text and produces transformed source by expanding generic type parameters, free-form macros, conditional blocks, and embedded HTML templates.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      # filesystem

codegen
  codegen.tokenize
    fn (source: string) -> result[list[token], string]
    + splits source into tokens preserving directive markers
    - returns error on unterminated string literals
    # lexing
  codegen.parse_directives
    fn (tokens: list[token]) -> result[list[directive], string]
    + extracts generic, macro, conditional, and template directives
    - returns error when a directive is malformed
    # parsing
  codegen.define_macro
    fn (state: codegen_state, name: string, params: list[string], body: string) -> codegen_state
    + registers a macro definition
    # macros
  codegen.expand_macro
    fn (state: codegen_state, name: string, args: list[string]) -> result[string, string]
    + substitutes arguments into a macro body
    - returns error when the macro is unknown
    - returns error when the argument count does not match
    # macros
  codegen.instantiate_generic
    fn (state: codegen_state, template_name: string, type_args: list[string]) -> result[string, string]
    + produces a monomorphized copy of a generic template for the given type args
    - returns error when the template is unknown
    # generics
  codegen.set_condition
    fn (state: codegen_state, flag: string, value: bool) -> codegen_state
    + sets a named flag used by conditional compilation
    # conditionals
  codegen.eval_condition
    fn (state: codegen_state, expr: string) -> result[bool, string]
    + evaluates a boolean expression over the defined flags
    - returns error when the expression references an undefined flag
    # conditionals
  codegen.render_template
    fn (state: codegen_state, template: string, data: map[string, string]) -> result[string, string]
    + interpolates a template with the given named values
    - returns error when a referenced placeholder has no value
    # templating
  codegen.transform
    fn (state: codegen_state, source: string) -> result[string, string]
    + runs the full pipeline: tokenize, parse, expand, emit
    - returns error at the first failing stage
    # pipeline
  codegen.transform_file
    fn (state: codegen_state, in_path: string, out_path: string) -> result[void, string]
    + reads a source file, transforms it, and writes the result
    - returns error when reading or writing fails
    # pipeline
    -> std.fs.read_all
    -> std.fs.write_all
