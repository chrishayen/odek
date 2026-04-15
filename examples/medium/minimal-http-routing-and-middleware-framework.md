# Requirement: "a minimalist HTTP routing and middleware framework"

Router with method-scoped routes, parameter paths, and a middleware chain.

std: (all units exist)

web
  web.new_router
    fn () -> router_state
    + creates a router with no routes and no middleware
    # construction
  web.use
    fn (r: router_state, middleware: fn(request, next_fn) -> response) -> router_state
    + appends a middleware to the global chain
    # middleware
  web.route
    fn (r: router_state, method: string, pattern: string, handler: fn(request) -> response) -> router_state
    + registers a handler for method and pattern; pattern supports ":name" segments
    ? supported methods include GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
    # registration
  web.group
    fn (r: router_state, prefix: string, build: fn(router_state) -> router_state) -> router_state
    + applies build inside a sub-router whose routes are prefixed
    # grouping
  web.match
    fn (r: router_state, method: string, path: string) -> optional[matched_route]
    + returns the matching route with extracted path parameters, or none
    # matching
  web.handle
    fn (r: router_state, req: request) -> response
    + runs the middleware chain followed by the matched handler
    + returns a 404 response when no route matches
    - returns a 405 response when the path matches but the method does not
    # dispatch
