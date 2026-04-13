# Requirement: "an HTTP request multiplexer"

Pattern-based router that dispatches incoming request method+path to a registered handler, extracting path parameters from patterns like `/users/:id`.

std: (all units exist)

mux
  mux.new
    @ () -> mux_state
    + creates an empty multiplexer with no routes
    # construction
  mux.register
    @ (state: mux_state, method: string, pattern: string, handler_id: string) -> result[mux_state, string]
    + adds a route; patterns may contain `:name` segments
    - returns error when the same method and pattern are already registered
    # registration
  mux.match
    @ (state: mux_state, method: string, path: string) -> optional[route_match]
    + returns the handler id and extracted params when the path matches a route
    - returns none when no route matches the method and path
    # dispatch
  mux.params_of
    @ (m: route_match) -> map[string, string]
    + returns the extracted path parameters for a match
    # inspection
  mux.handler_of
    @ (m: route_match) -> string
    + returns the handler id for a match
    # inspection
