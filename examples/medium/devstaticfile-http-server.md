# Requirement: "a developer-oriented local static-file HTTP server"

Provides route registration and request dispatch; a host transport layer supplies the socket.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file at the given path
      - returns error when the file is missing or unreadable
      # filesystem
  std.path
    std.path.join
      fn (base: string, rel: string) -> string
      + joins two path segments with a platform-appropriate separator
      # paths
    std.path.is_within
      fn (root: string, candidate: string) -> bool
      + returns true when candidate resolves inside root
      ? used to reject directory traversal attempts
      # paths

devserver
  devserver.new
    fn (root_dir: string) -> server_state
    + creates a server rooted at the given directory with no custom routes
    # construction
  devserver.add_route
    fn (state: server_state, path: string, body: bytes, content_type: string) -> server_state
    + registers a fixed response for the given request path
    # routing
  devserver.handle
    fn (state: server_state, method: string, request_path: string) -> http_response
    + returns a static-file response by joining request_path to root
    + returns a registered route's response when one matches
    - returns 404 when no route matches and the file does not exist
    - returns 403 when the resolved path escapes the root directory
    - returns 405 when method is anything other than "GET" or "HEAD"
    # dispatch
    -> std.path.join
    -> std.path.is_within
    -> std.fs.read_all
  devserver.response_bytes
    fn (response: http_response) -> bytes
    + serializes the response into wire-format HTTP/1.1 bytes
    # serialization
