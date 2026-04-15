# Requirement: "a minimalist web framework"

A router-centric web framework: method+path routes, parameter extraction, ordered middleware, request context, and a response builder. The HTTP transport itself is the caller's concern; this library produces handler dispatch over a parsed request/response pair.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits a string on the literal separator
      # strings
    std.strings.starts_with
      fn (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
  std.json
    std.json.encode
      fn (value: string) -> string
      + encodes a string as a JSON document
      # serialization
    std.json.parse
      fn (raw: string) -> result[string, string]
      + parses a JSON document into its textual form
      - returns error on invalid JSON
      # serialization
  std.url
    std.url.parse_query
      fn (raw: string) -> map[string, string]
      + parses a url-encoded query string into a map
      # url

framework
  framework.new
    fn () -> framework_state
    + creates an empty router with no routes or middleware
    # construction
  framework.get
    fn (state: framework_state, path_pattern: string, handler_id: string) -> framework_state
    + registers a handler for GET on the pattern
    # route_registration
  framework.post
    fn (state: framework_state, path_pattern: string, handler_id: string) -> framework_state
    + registers a handler for POST on the pattern
    # route_registration
  framework.put
    fn (state: framework_state, path_pattern: string, handler_id: string) -> framework_state
    + registers a handler for PUT on the pattern
    # route_registration
  framework.delete
    fn (state: framework_state, path_pattern: string, handler_id: string) -> framework_state
    + registers a handler for DELETE on the pattern
    # route_registration
  framework.use
    fn (state: framework_state, middleware_id: string) -> framework_state
    + appends middleware to the ordered chain
    # middleware
  framework.group
    fn (state: framework_state, prefix: string) -> framework_state
    + creates a scoped subrouter sharing middleware
    # grouping
  framework.match
    fn (state: framework_state, method: string, path: string) -> result[route_match, string]
    + returns handler id and extracted path params on match
    - returns error when no route matches
    # routing
    -> std.strings.split
    -> std.strings.starts_with
  framework.build_context
    fn (method: string, path: string, query: string, body: string) -> request_context
    + builds the per-request context passed to handlers
    # context
    -> std.url.parse_query
  framework.param
    fn (ctx: request_context, name: string) -> optional[string]
    + returns a path or query parameter by name
    # context
  framework.bind_json
    fn (ctx: request_context) -> result[string, string]
    + parses request body as JSON and returns the raw value
    - returns error on malformed JSON
    # binding
    -> std.json.parse
  framework.respond_json
    fn (ctx: request_context, status: i32, body: string) -> response
    + builds a response with status code and JSON body
    # response
    -> std.json.encode
  framework.respond_text
    fn (ctx: request_context, status: i32, body: string) -> response
    + builds a response with status code and plain-text body
    # response
  framework.dispatch
    fn (state: framework_state, ctx: request_context) -> response
    + runs middleware chain then matched handler and returns the response
    - produces 404 response when no route matches
    # dispatch
