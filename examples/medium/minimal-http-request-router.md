# Requirement: "a minimal HTTP request router"

Routes bind method+path patterns to handler ids. The router parses a raw request line and returns the matched handler along with extracted path parameters. The caller provides the transport and handlers.

std
  std.text
    std.text.split_path_segments
      fn (path: string) -> list[string]
      + splits on '/' and drops empty segments
      # text
  std.http
    std.http.parse_request_line
      fn (line: string) -> result[tuple[string, string, string], string]
      + returns (method, path, version) from "GET /foo HTTP/1.1"
      - returns error when the line does not have three fields
      # http
    std.http.status_text
      fn (code: i32) -> string
      + returns the canonical reason phrase for a status code
      # http

http_router
  http_router.new
    fn () -> router_state
    + returns an empty router
    # construction
  http_router.add_route
    fn (state: router_state, method: string, pattern: string, handler_id: string) -> router_state
    + registers a handler for method+pattern
    ? segments of the form ":name" become path parameters
    # routing
  http_router.match
    fn (state: router_state, method: string, path: string) -> optional[tuple[string, map[string, string]]]
    + returns (handler_id, params) on match
    - returns empty when no route matches
    # routing
    -> std.text.split_path_segments
  http_router.dispatch_request_line
    fn (state: router_state, line: string) -> result[tuple[string, map[string, string]], string]
    + parses the request line and looks up the handler
    - returns error when the request line is malformed
    - returns error when no route matches
    # dispatch
    -> std.http.parse_request_line
  http_router.format_status_line
    fn (code: i32) -> string
    + returns "HTTP/1.1 <code> <reason>"
    # formatting
    -> std.http.status_text
