# Requirement: "a minimalistic HTTP request multiplexer with support for request-scoped context"

Holds a table of (method, pattern) to handler, dispatches on match, and threads a context map through to handlers.

std: (all units exist)

mux
  mux.new
    @ () -> mux_state
    + creates an empty multiplexer with no routes
    # construction
  mux.handle
    @ (m: mux_state, method: string, pattern: string, handler: fn(http_request, map[string, string]) -> http_response) -> mux_state
    + registers a handler for a (method, pattern) pair
    ? pattern supports "/prefix/:name" style parameters
    # routing
  mux.dispatch
    @ (m: mux_state, req: http_request, ctx: map[string, string]) -> http_response
    + finds the best-matching route and invokes its handler with the merged path parameters
    - returns a 404 response when no route matches
    - returns a 405 response when the path matches but the method does not
    # routing
  mux.match_pattern
    @ (pattern: string, path: string) -> optional[map[string, string]]
    + returns the captured parameters when path matches pattern
    - returns none when the pattern does not match
    # pattern_matching
