# Requirement: "a framework for building programming languages and language-server protocol backends"

A reusable compiler front-end plus an LSP event loop. The project layer exposes the main building blocks: source document management, lexing, parsing, symbol resolution, and diagnostic reporting. The LSP transport layer lives behind thin std primitives.

std
  std.io
    std.io.read_line
      @ () -> result[string, string]
      + reads one line from standard input and returns it without the trailing newline
      - returns error on end of input
      # io
    std.io.write_string
      @ (s: string) -> void
      + writes the string to standard output
      # io
  std.json
    std.json.parse_value
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_value
      @ (value: map[string, string]) -> string
      + encodes a string-to-string map as a JSON object
      # serialization

language_framework
  language_framework.new_document_store
    @ () -> document_store
    + returns an empty document store
    # construction
  language_framework.open_document
    @ (store: document_store, uri: string, text: string) -> document_store
    + registers the document under the given uri with its initial text
    # documents
  language_framework.update_document
    @ (store: document_store, uri: string, text: string) -> result[document_store, string]
    + replaces the document text and invalidates cached analysis
    - returns error when the uri is not open
    # documents
  language_framework.tokenize
    @ (text: string) -> list[token]
    + produces tokens with kind, start offset, and length
    + returns an empty list for empty input
    # lexing
  language_framework.parse
    @ (tokens: list[token]) -> result[ast_node, list[diagnostic]]
    + returns the root ast node when the token stream is well-formed
    - returns diagnostics describing syntax errors with source offsets
    # parsing
  language_framework.resolve_symbols
    @ (root: ast_node) -> symbol_table
    + walks the ast and records every declared name with its scope
    # analysis
  language_framework.find_definition
    @ (table: symbol_table, name: string, offset: i32) -> optional[source_range]
    + returns the declaration range for the symbol at the given offset
    - returns none when the symbol is undefined at that position
    # analysis
  language_framework.collect_diagnostics
    @ (store: document_store, uri: string) -> list[diagnostic]
    + runs lexing, parsing, and resolution and returns all diagnostics
    - returns an empty list when the uri is not open
    # diagnostics
  language_framework.new_server
    @ (store: document_store) -> lsp_server_state
    + returns an initialized language server bound to the store
    # construction
  language_framework.handle_message
    @ (server: lsp_server_state, raw: string) -> tuple[lsp_server_state, optional[string]]
    + dispatches an LSP JSON-RPC message and returns any response to send back
    - returns (unchanged_server, none) when the message is a notification with no reply
    # lsp
    -> std.json.parse_value
    -> std.json.encode_value
  language_framework.run_stdio_loop
    @ (server: lsp_server_state) -> result[void, string]
    + reads LSP framed messages from stdin, dispatches them, and writes replies to stdout
    - returns error when the input stream breaks mid-message
    # transport
    -> std.io.read_line
    -> std.io.write_string
