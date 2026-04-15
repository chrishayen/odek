# Requirement: "an http router with path parameters and middleware"

A router that matches method and path against registered routes, extracts path parameters, and runs a middleware chain before the handler.

std: (all units exist)

http_router
  http_router.new
    fn () -> router_state
    + returns an empty router
    # construction
  http_router.add_route
    fn (router: router_state, method: string, pattern: string, handler_id: string) -> router_state
    + returns a router with the route added
    ? patterns use "/users/:id" style placeholders
    # registration
  http_router.use
    fn (router: router_state, middleware_id: string) -> router_state
    + appends a middleware to the global chain
    # middleware
  http_router.group
    fn (router: router_state, prefix: string, middleware_ids: list[string]) -> router_state
    + returns a router where subsequent routes inherit the prefix and middleware
    # grouping
  http_router.match
    fn (router: router_state, method: string, path: string) -> optional[route_match]
    + returns the matched handler id and extracted parameters
    - returns none when no route matches
    - returns none when the method does not match
    # dispatch
  http_router.params_of
    fn (match: route_match) -> map[string, string]
    + returns the parameters extracted from the matched path
    # parameters
  http_router.middleware_chain
    fn (router: router_state, match: route_match) -> list[string]
    + returns the ordered list of middleware ids that apply to the match
    # middleware
