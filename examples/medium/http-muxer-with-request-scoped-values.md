# Requirement: "a high performance HTTP multiplexer with request-scoped values"

A trie-based router that matches method and path, extracts path parameters, and exposes a request-scoped value map for handlers.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits a string on the literal separator
      # strings

mux
  mux.new
    fn () -> mux_state
    + creates an empty multiplexer
    # construction
  mux.handle
    fn (state: mux_state, method: string, pattern: string, handler_id: string) -> mux_state
    + registers a handler for the (method, pattern) pair
    ? pattern segments starting with ":" are parameter placeholders
    # registration
    -> std.strings.split
  mux.lookup
    fn (state: mux_state, method: string, path: string) -> result[route_match, string]
    + returns the handler id and path parameters on match
    - returns error when method or path has no registered handler
    # routing
    -> std.strings.split
  mux.new_context
    fn () -> request_context
    + creates an empty request-scoped value map
    # context
  mux.context_set
    fn (ctx: request_context, key: string, value: string) -> request_context
    + stores a value under a key in the request context
    # context
  mux.context_get
    fn (ctx: request_context, key: string) -> optional[string]
    + retrieves a value by key
    - returns none when key is absent
    # context
