# Requirement: "an HTTP request multiplexer that dispatches requests to handlers by method and path"

Registers handlers against a method and path pattern, then matches incoming requests. std contributes only a path splitter; the mux owns routing logic.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s by the separator, returning non-separator segments
      + empty input returns a list containing one empty string
      # strings

mux
  mux.new
    fn () -> mux_state
    + creates an empty multiplexer with no registered routes
    # construction
  mux.handle
    fn (m: mux_state, method: string, pattern: string, handler_id: string) -> result[void, string]
    + registers a handler for the method+pattern pair
    - returns error when the method+pattern is already registered
    - returns error when pattern does not start with "/"
    # registration
    -> std.strings.split
  mux.match
    fn (m: mux_state, method: string, path: string) -> optional[match_result]
    + returns the matched handler id and extracted path params for method+path
    - returns none when no registered route matches
    ? literal segments beat ":param" segments beat "*rest" wildcards
    # routing
    -> std.strings.split
  mux.routes
    fn (m: mux_state) -> list[route_info]
    + returns all registered routes as (method, pattern, handler_id) tuples
    # inspection
