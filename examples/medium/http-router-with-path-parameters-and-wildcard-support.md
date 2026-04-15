# Requirement: "an HTTP router with path parameters and wildcard support"

A router that compiles registered patterns into a trie for fast longest-prefix matching. Only one thin std helper is needed for path splitting.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s by the separator
      + returns one empty segment for empty input
      # strings

router
  router.new
    fn () -> router_state
    + creates an empty router with no routes
    # construction
  router.add
    fn (r: router_state, method: string, pattern: string, handler_id: string) -> result[void, string]
    + registers a handler under method+pattern, inserting into the trie
    - returns error when the exact method+pattern is already registered
    - returns error when pattern does not start with "/"
    # registration
    -> std.strings.split
  router.lookup
    fn (r: router_state, method: string, path: string) -> optional[lookup_result]
    + returns the matched handler id and captured parameters
    - returns none when no route matches for the method
    ? static segments > ":name" params > "*rest" wildcards
    # routing
    -> std.strings.split
  router.allowed_methods
    fn (r: router_state, path: string) -> list[string]
    + returns the methods registered for the path (useful for 405 responses)
    # inspection
