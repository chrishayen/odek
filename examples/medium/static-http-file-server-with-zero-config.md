# Requirement: "a zero-config static HTTP file server"

Maps an incoming request path to a file under a root directory and builds the response metadata. Actual socket I/O is the caller's responsibility; the library answers "what should the response be?"

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file at path
      - returns error when the path does not exist or is unreadable
      # filesystem
    std.fs.is_dir
      fn (path: string) -> bool
      + returns true when path exists and is a directory
      # filesystem
  std.path
    std.path.join
      fn (a: string, b: string) -> string
      + joins two path segments with the platform separator, collapsing duplicates
      # paths
    std.path.extension
      fn (p: string) -> string
      + returns the lowercase extension including the leading dot, or "" if none
      # paths

static_server
  static_server.new
    fn (root: string) -> server_state
    + creates a server rooted at the given directory
    # construction
  static_server.resolve
    fn (state: server_state, request_path: string) -> result[string, string]
    + joins request_path onto the root and returns the resolved file path
    - returns error when request_path escapes the root via `..`
    # routing
    -> std.path.join
  static_server.mime_for
    fn (path: string) -> string
    + returns a mime type based on the file extension, defaulting to "application/octet-stream"
    # content_negotiation
    -> std.path.extension
  static_server.serve
    fn (state: server_state, request_path: string) -> result[http_response, http_response]
    + returns a 200 response with file bytes and mime type on success
    - returns a 403 when the path escapes the root
    - returns a 404 when the file does not exist
    - returns a 405-style directory listing refusal when path is a directory
    # response_building
    -> static_server.resolve
    -> static_server.mime_for
    -> std.fs.read_all
    -> std.fs.is_dir
