# Requirement: "a static file HTTP server"

A server that maps URL paths to files under a root directory and streams them with correct content types.

std
  std.net
    std.net.http_listen
      fn (host: string, port: i32) -> result[listener_handle, string]
      + binds an HTTP listener
      - returns error when the port is in use
      # networking
    std.net.http_accept
      fn (lis: listener_handle) -> result[http_request, string]
      + blocks until a request arrives
      # networking
    std.net.http_respond
      fn (req: http_request, status: i32, headers: map[string, string], body: bytes) -> result[void, string]
      + writes a response
      # networking
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns directory entry names
      - returns error when the directory does not exist
      # filesystem
    std.fs.is_dir
      fn (path: string) -> bool
      + returns true when the path exists and is a directory
      # filesystem
  std.path
    std.path.join
      fn (base: string, rel: string) -> string
      + joins two path segments
      # paths
    std.path.is_safe_subpath
      fn (base: string, candidate: string) -> bool
      + returns true when candidate resolves inside base after normalization
      # paths

file_server
  file_server.new
    fn (root_dir: string) -> server_state
    + creates a server rooted at a directory
    # construction
  file_server.resolve
    fn (srv: server_state, url_path: string) -> result[string, string]
    + maps a URL path to an absolute filesystem path
    - returns error when the resolved path escapes the root
    # resolution
    -> std.path.join
    -> std.path.is_safe_subpath
  file_server.content_type
    fn (path: string) -> string
    + returns a MIME type based on the file extension, defaulting to application/octet-stream
    # mime
  file_server.render_index
    fn (dir: string, entries: list[string]) -> bytes
    + renders a simple HTML directory listing
    # rendering
  file_server.handle_request
    fn (srv: server_state, req: http_request) -> result[void, string]
    + resolves the path, reads the file or lists the directory, and writes the response
    - returns error when reading the file fails
    # serving
    -> std.fs.read_all
    -> std.fs.is_dir
    -> std.fs.list_dir
    -> std.net.http_respond
  file_server.listen
    fn (srv: server_state, host: string, port: i32) -> result[void, string]
    + accepts requests and serves them from the root directory
    - returns error when the listener cannot be bound
    # serving
    -> std.net.http_listen
    -> std.net.http_accept
