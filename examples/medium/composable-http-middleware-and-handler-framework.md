# Requirement: "a composable http middleware and handler framework"

A small framework for chaining middleware around request handlers and dispatching by method and path.

std: (all units exist)

middleware_chain
  middleware_chain.new_router
    fn () -> router_state
    + creates an empty router with no routes
    # construction
  middleware_chain.add_route
    fn (state: router_state, method: string, path: string, handler: handler_fn) -> router_state
    + registers a handler for an exact method and path
    # routing
  middleware_chain.use
    fn (state: router_state, mw: middleware_fn) -> router_state
    + appends a middleware to the global chain
    ? middleware wrap every handler in the order added
    # middleware
  middleware_chain.dispatch
    fn (state: router_state, req: request_obj) -> response_obj
    + runs global middleware then the matched handler
    - returns a 404 response when no route matches
    - returns a 405 response when the path matches but the method does not
    # dispatch
  middleware_chain.compose
    fn (outer: middleware_fn, inner: middleware_fn) -> middleware_fn
    + returns a single middleware that runs outer before inner
    # composition
