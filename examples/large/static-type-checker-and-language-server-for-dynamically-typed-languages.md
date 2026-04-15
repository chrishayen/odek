# Requirement: "a static type checker and language server for a dynamically typed language"

A checker pipeline (lex, parse, resolve, infer, diagnose) plus a minimal language-server protocol surface. Std provides generic IO and JSON-RPC primitives.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns file contents as a string
      - returns error when the file does not exist
      # filesystem
  std.io
    std.io.read_line
      fn () -> result[string, string]
      + reads a line from stdin
      # io
    std.io.write_all
      fn (data: string) -> void
      + writes data to stdout
      # io
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document
      - returns error on malformed input
      # serialization
    std.json.encode
      fn (value: json_value) -> string
      + encodes a JSON value as text
      # serialization

checker
  checker.tokenize
    fn (source: string) -> result[list[token], string]
    + splits source into tokens with line and column info
    - returns error on unterminated strings
    # lexing
  checker.parse
    fn (tokens: list[token]) -> result[ast_module, string]
    + builds a module AST from tokens
    - returns error with location on unexpected tokens
    # parsing
  checker.resolve_names
    fn (module: ast_module) -> result[resolved_module, list[diagnostic]]
    + binds each name reference to its definition
    - reports undefined name diagnostics
    # name_resolution
  checker.infer_types
    fn (module: resolved_module) -> tuple[typed_module, list[diagnostic]]
    + infers a type for every expression using Hindley-Milner style unification
    + returns any constraint failures as diagnostics
    # type_inference
  checker.check_file
    fn (path: string) -> result[list[diagnostic], string]
    + runs the full pipeline on one file and returns all diagnostics
    - returns error when the file cannot be read
    # checking
    -> std.fs.read_all
  checker.check_project
    fn (root: string, files: list[string]) -> map[string, list[diagnostic]]
    + checks a set of files and returns diagnostics keyed by path
    ? files are checked independently; cross-file types are resolved via resolve_names
    # checking

language_server
  language_server.new
    fn (workspace_root: string) -> server_state
    + creates a server bound to a workspace root
    # construction
  language_server.handle_initialize
    fn (state: server_state, params: map[string, string]) -> tuple[server_state, string]
    + returns server capabilities as a JSON string
    # lsp_initialize
    -> std.json.encode
  language_server.handle_did_open
    fn (state: server_state, uri: string, text: string) -> tuple[server_state, list[diagnostic]]
    + stores the document and returns fresh diagnostics
    # lsp_textdocument
  language_server.handle_did_change
    fn (state: server_state, uri: string, text: string) -> tuple[server_state, list[diagnostic]]
    + replaces document contents and rechecks
    # lsp_textdocument
  language_server.handle_hover
    fn (state: server_state, uri: string, line: i32, column: i32) -> optional[string]
    + returns the inferred type at a position, if any
    - returns none when the position is outside any expression
    # lsp_hover
  language_server.serve_loop
    fn (state: server_state) -> result[void, string]
    + reads framed JSON-RPC messages from stdin and dispatches them
    - returns error on malformed framing
    # lsp_transport
    -> std.io.read_line
    -> std.io.write_all
    -> std.json.parse
    -> std.json.encode
