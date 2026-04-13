# Requirement: "a high performance HTTP router"

A trie-based method+path router with parameter capture and wildcard segments. Handlers are referenced by id; the caller binds ids to functions at its own layer.

std
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits a string on the literal separator
      # strings

router
  router.new
    @ () -> router_state
    + creates an empty router
    # construction
  router.add
    @ (state: router_state, method: string, pattern: string, handler_id: string) -> result[router_state, string]
    + registers a handler for the (method, pattern) pair
    - returns error on duplicate registration
    ? segments starting with ":" are parameters, "*" matches trailing path
    # registration
    -> std.strings.split
  router.lookup
    @ (state: router_state, method: string, path: string) -> result[route_match, string]
    + returns handler id and captured parameters on match
    - returns error when no route matches
    # routing
    -> std.strings.split
  router.allowed_methods
    @ (state: router_state, path: string) -> list[string]
    + returns HTTP methods registered for the path
    - returns empty list when the path has no routes
    # introspection
    -> std.strings.split
