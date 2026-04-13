# Requirement: "an API gateway framework that routes incoming HTTP requests to one or more backend endpoints through a configurable middleware chain"

Routes are declared with a path pattern, a list of backends, and a middleware stack. The gateway matches a request, runs middleware, fans out to backends, and merges their responses.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body
      - returns error on malformed request line or headers
      # http
    std.http.serialize_response
      @ (status: i32, headers: map[string,string], body: bytes) -> bytes
      + returns a wire-format response
      # http
    std.http.client_send
      @ (host: string, port: i32, req: http_request) -> result[http_response, string]
      + performs an outbound request and returns the parsed response
      - returns error on network or protocol failure
      # http
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses JSON
      - returns error on invalid JSON
      # parsing
    std.json.encode
      @ (value: json_value) -> string
      # serialization

gateway
  gateway.new_router
    @ () -> router_state
    + returns an empty router
    # construction
  gateway.add_route
    @ (state: router_state, method: string, path_pattern: string, backends: list[backend_spec], middleware: list[string]) -> router_state
    + registers a route with its backends and middleware identifiers
    # registration
  gateway.match_route
    @ (state: router_state, method: string, path: string) -> optional[matched_route]
    + returns the first route whose method and path pattern match
    ? patterns support "/:param" segments and "/*rest" suffixes
    # routing
  gateway.extract_params
    @ (pattern: string, path: string) -> map[string,string]
    + returns the captured path parameters
    # routing
  gateway.run_middleware
    @ (chain: list[string], req: http_request, handlers: map[string, fn(http_request) -> result[http_request,http_response]]) -> result[http_request, http_response]
    + runs each middleware in order, short-circuiting on an early response
    ? middleware may mutate the request or return a terminal response
    # middleware
  gateway.fanout_backends
    @ (backends: list[backend_spec], req: http_request) -> result[list[http_response], string]
    + sends the request to every backend and collects responses
    - returns error when all backends fail
    # fanout
    -> std.http.client_send
  gateway.merge_responses
    @ (responses: list[http_response], strategy: string) -> http_response
    + merges fan-out responses per strategy ("first", "concat_json", "merge_json")
    ? unknown strategies fall back to "first"
    # merging
    -> std.json.parse
    -> std.json.encode
  gateway.handle
    @ (state: router_state, raw: bytes, handlers: map[string, fn(http_request) -> result[http_request,http_response]]) -> bytes
    + parses, routes, runs middleware, fans out, merges, and serializes the response
    + returns a 404 response when no route matches
    # pipeline
    -> std.http.parse_request
    -> std.http.serialize_response
    -> gateway.match_route
    -> gateway.run_middleware
    -> gateway.fanout_backends
    -> gateway.merge_responses
