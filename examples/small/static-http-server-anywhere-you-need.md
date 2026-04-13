# Requirement: "a static http server"

Serves files from a root directory over HTTP. The project layer maps request paths to files and builds responses; the network loop and filesystem primitives live in std.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.path_join
      @ (base: string, child: string) -> string
      + joins two path segments with the platform separator
      + collapses redundant separators
      # filesystem
  std.http
    std.http.serve
      @ (addr: string, handler: request_handler) -> result[void, string]
      + binds the address and dispatches each request to the handler
      - returns error when the address cannot be bound
      # network

static_server
  static_server.resolve_path
    @ (root: string, url_path: string) -> result[string, string]
    + joins the URL path onto the root directory
    - returns error when the resolved path escapes the root
    # path_safety
    -> std.fs.path_join
  static_server.content_type_for
    @ (path: string) -> string
    + returns "text/html" for .html, "text/css" for .css, and so on
    + returns "application/octet-stream" for unknown extensions
    # mime
  static_server.handle
    @ (root: string, url_path: string) -> result[tuple[i32, string, bytes], string]
    + returns (200, content_type, body) for an existing file
    - returns (404, "text/plain", "not found") when the file is missing
    - returns (403, "text/plain", "forbidden") when the path escapes the root
    # request_handling
    -> std.fs.read_all
