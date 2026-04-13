# Requirement: "a zero-configuration static file HTTP server"

Given a root directory, serves files over HTTP. Maps URL paths to files, guesses content types from extensions, and prevents path traversal.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file is missing
      # filesystem
    std.fs.stat
      @ (path: string) -> result[file_info, string]
      + returns size, mtime, and is_dir for the path
      - returns error when the path does not exist
      # filesystem
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses an HTTP/1.1 request
      - returns error on malformed input
      # http_parsing
    std.http.encode_response
      @ (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes an HTTP/1.1 response
      # http_encoding
  std.net
    std.net.listen_tcp
      @ (host: string, port: i32) -> result[listener, string]
      + binds and listens on the given address
      # networking
    std.net.accept
      @ (l: listener) -> result[connection, string]
      + blocks until a connection arrives
      # networking

static_server
  static_server.resolve_path
    @ (root: string, url_path: string) -> result[string, string]
    + returns root joined with url_path, canonicalized
    - returns error when the resolved path escapes root
    # path_resolution
  static_server.content_type_for
    @ (path: string) -> string
    + returns a MIME type based on the file extension
    ? unknown extensions return "application/octet-stream"
    # content_type
  static_server.handle
    @ (root: string, req: http_request) -> http_response
    + returns 200 with the file body when the path resolves to a regular file
    - returns 404 when the file does not exist
    - returns 403 when the path would escape root
    # request_handling
    -> std.fs.read_all
    -> std.fs.stat
  static_server.serve
    @ (root: string, host: string, port: i32) -> result[void, string]
    + accepts connections and serves each request from root
    - returns error when the listener cannot bind
    # serving
    -> std.net.listen_tcp
    -> std.net.accept
    -> std.http.parse_request
    -> std.http.encode_response
