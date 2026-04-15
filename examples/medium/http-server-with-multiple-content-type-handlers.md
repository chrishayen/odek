# Requirement: "an http web server with multiple content-type handlers"

A minimal router that dispatches by path and content-type handler. Markdown and a css-like language compile to html on request.

std
  std.io
    std.io.read_file
      fn (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the path does not exist
      # filesystem
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses method, path, headers from a raw request
      - returns error on malformed request line
      # http
    std.http.write_response
      fn (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes a response with the given status, headers, and body
      # http

web_server
  web_server.new_router
    fn () -> router_state
    + returns an empty router with no routes registered
    # construction
  web_server.register
    fn (state: router_state, path_prefix: string, handler_kind: string) -> router_state
    + binds a prefix to one of: "static", "markdown", "stylesheet"
    - returns unchanged state when handler_kind is unknown
    # routing
  web_server.handle
    fn (state: router_state, req: http_request) -> bytes
    + dispatches to the handler bound to the longest matching prefix
    + returns a 404 response when no route matches
    # dispatch
    -> std.http.write_response
    -> std.io.read_file
  web_server.render_markdown
    fn (source: string) -> string
    + converts headings, paragraphs, and emphasis to html
    + preserves plain text outside markup
    # markdown
  web_server.render_stylesheet
    fn (source: string) -> string
    + translates nested selector blocks into flat css rules
    - skips blocks whose selector is empty
    # stylesheet
