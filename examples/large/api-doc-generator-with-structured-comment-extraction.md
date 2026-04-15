# Requirement: "an API documentation generator that extracts structured comments from source files"

Scans a directory of source files, extracts doc comments attached to declarations, builds a model, and renders HTML.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns all regular file paths under root
      - returns error when root does not exist
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the full contents of the file as a string
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes the full contents to the file, creating directories as needed
      # filesystem
  std.html
    std.html.escape
      fn (raw: string) -> string
      + escapes &, <, >, ", and ' for safe HTML embedding
      # html

docgen
  docgen.tokenize_source
    fn (source: string) -> list[doc_token]
    + returns a stream of tokens including doc comments, declarations, and braces
    ? doc comments are delimited by /** */
    # lexing
  docgen.extract_comments
    fn (tokens: list[doc_token]) -> list[raw_comment]
    + pairs each doc comment with the declaration that follows it
    # extraction
  docgen.parse_comment
    fn (raw: raw_comment) -> doc_entry
    + splits the comment into a summary, description, and tag block
    + recognizes @param, @returns, @throws, @example
    # parsing
  docgen.parse_signature
    fn (raw: raw_comment) -> doc_signature
    + extracts the declaration name, kind, and parameter list from the following declaration
    # parsing
  docgen.build_model
    fn (root: string) -> result[doc_model, string]
    + walks the directory, parses every file, and returns a grouped model by module
    - returns error when no source files are found
    # model
    -> std.fs.walk
    -> std.fs.read_all
  docgen.render_entry
    fn (entry: doc_entry, sig: doc_signature) -> string
    + renders a single entry as an HTML fragment
    # rendering
    -> std.html.escape
  docgen.render_index
    fn (model: doc_model) -> string
    + renders the module index page with links to each entry
    # rendering
    -> std.html.escape
  docgen.render_module
    fn (model: doc_model, module: string) -> string
    + renders the page for a single module containing all its entries
    # rendering
    -> std.html.escape
  docgen.generate
    fn (source_root: string, out_dir: string) -> result[i32, string]
    + builds the model and writes index and per-module HTML files, returning the number of pages written
    - returns error when the output directory cannot be created
    # pipeline
    -> std.fs.write_all
