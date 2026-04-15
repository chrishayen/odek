# Requirement: "an API documentation generator with support for modern class syntax and type annotations"

Similar to a generic doc generator, but the parser understands type annotations on parameters and return types and carries them through to the rendered output.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns all regular file paths under root
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the full file contents
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes the full contents to path
      # filesystem
  std.html
    std.html.escape
      fn (raw: string) -> string
      + escapes HTML metacharacters
      # html

typeddoc
  typeddoc.tokenize
    fn (source: string) -> list[doc_token]
    + produces a token stream including class, method, comment, and type annotation tokens
    # lexing
  typeddoc.parse_type_annotation
    fn (raw: string) -> type_expr
    + parses generics, unions, arrays, and record types into a structured expression
    ? unrecognized shapes become an opaque "unknown" type
    # type_parsing
  typeddoc.extract_declarations
    fn (tokens: list[doc_token]) -> list[typed_decl]
    + returns every class, method, and top-level function paired with its preceding doc comment
    # extraction
  typeddoc.parse_doc_tags
    fn (comment: string) -> map[string, string]
    + parses leading-@ tag blocks into a tag-to-body map
    # parsing
  typeddoc.merge_annotations
    fn (decl: typed_decl, tags: map[string,string]) -> typed_entry
    + overlays inline type annotations with explicit @param and @returns types from tags
    + the inline annotation wins on conflict
    # merging
  typeddoc.build_model
    fn (root: string) -> result[typed_model, string]
    + walks the source tree and assembles the full typed model
    - returns error when no sources are found
    # model
    -> std.fs.walk
    -> std.fs.read_all
  typeddoc.format_type
    fn (expr: type_expr) -> string
    + renders a type expression back to a canonical display string
    # rendering
  typeddoc.render_class
    fn (cls: typed_entry) -> string
    + renders a class page including all methods and their signatures
    # rendering
    -> std.html.escape
  typeddoc.render_method
    fn (method: typed_entry) -> string
    + renders a method entry with parameter and return types
    # rendering
    -> std.html.escape
  typeddoc.generate
    fn (source_root: string, out_dir: string) -> result[i32, string]
    + builds the model and writes all HTML pages, returning the file count
    - returns error when the output directory cannot be created
    # pipeline
    -> std.fs.write_all
