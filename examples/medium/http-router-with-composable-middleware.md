# Requirement: "an HTTP router with composable middleware"

A trie-based router with method dispatch and an ordered middleware chain. Handlers receive a context bag so dependency wiring is just map lookups.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits a string on the separator
      + returns a single-element list when separator is absent
      # strings

router
  router.new
    fn () -> router_state
    + returns an empty router with no registered routes
    # construction
  router.handle
    fn (state: router_state, method: string, pattern: string, handler: handler_fn) -> router_state
    + registers a handler for an exact method and path pattern
    + supports path parameters like "/users/:id"
    # registration
    -> std.strings.split
  router.use
    fn (state: router_state, mw: middleware_fn) -> router_state
    + appends middleware to the chain in registration order
    # middleware
  router.group
    fn (state: router_state, prefix: string, build: group_builder_fn) -> router_state
    + registers a subtree of routes sharing a common prefix and middleware
    # grouping
  router.lookup
    fn (state: router_state, method: string, path: string) -> optional[matched_route]
    + returns the matched handler and extracted params when a route matches
    - returns none when no route matches the method and path
    # lookup
    -> std.strings.split
  router.dispatch
    fn (state: router_state, method: string, path: string, ctx: request_ctx) -> response
    + runs the middleware chain and invokes the matched handler
    - returns a 404 response when no route matches
    - returns a 405 response when the path matches but method does not
    # dispatch
