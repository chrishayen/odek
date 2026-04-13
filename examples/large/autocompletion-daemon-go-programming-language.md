# Requirement: "a source code autocompletion daemon"

Parses source files, maintains a symbol index, serves completion queries over a local socket.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file as text
      - returns error when missing
      # io
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns every file path under root
      # io
  std.net
    std.net.unix_listen
      @ (path: string) -> result[listener_state, string]
      + binds a Unix domain socket for incoming requests
      - returns error when the path is already in use
      # networking
    std.net.accept
      @ (listener: listener_state) -> result[conn_state, string]
      + blocks for the next incoming connection
      # networking
    std.net.conn_read_line
      @ (conn: conn_state) -> result[string, string]
      + reads a newline-delimited request
      # networking
    std.net.conn_write
      @ (conn: conn_state, data: bytes) -> result[i32, string]
      + writes a response
      # networking
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string map
      - returns error on malformed input
      # serialization
    std.json.encode_array
      @ (items: list[string]) -> string
      + encodes a list of strings as a JSON array
      # serialization

gocode
  gocode.index_new
    @ () -> symbol_index
    + creates an empty symbol index
    # construction
  gocode.tokenize
    @ (source: string) -> list[token]
    + splits source text into identifier and punctuation tokens
    # lexing
  gocode.parse_declarations
    @ (tokens: list[token]) -> list[symbol]
    + extracts function, type, and variable declarations
    # parsing
  gocode.index_add_file
    @ (index: symbol_index, path: string) -> result[symbol_index, string]
    + parses a file and records its symbols
    - returns error when the file cannot be read
    # indexing
    -> std.fs.read_all
  gocode.index_project
    @ (index: symbol_index, root: string) -> result[symbol_index, string]
    + walks a project directory and indexes every source file
    # indexing
    -> std.fs.walk
  gocode.query
    @ (index: symbol_index, prefix: string, scope: string) -> list[completion]
    + returns symbols whose names start with the prefix, scored by proximity to scope
    # querying
  gocode.serve
    @ (index: symbol_index, socket_path: string) -> result[void, string]
    + accepts completion requests on a Unix socket and writes JSON responses
    - returns error when the listener fails
    # serving
    -> std.net.unix_listen
    -> std.net.accept
    -> std.net.conn_read_line
    -> std.net.conn_write
    -> std.json.parse_object
    -> std.json.encode_array
  gocode.shutdown
    @ (index: symbol_index) -> void
    + releases resources held by the index
    # lifecycle
