# Requirement: "a minimal HTTP router"

Matches method and exact path to a handler identifier. Keeps the surface intentionally small.

std: (all units exist)

http_router
  http_router.new
    fn () -> router_state
    + returns an empty router
    # construction
  http_router.handle
    fn (state: router_state, method: string, path: string, handler_id: string) -> router_state
    + registers a handler for the method and exact path
    # registration
  http_router.lookup
    fn (state: router_state, method: string, path: string) -> optional[string]
    + returns the handler id registered for the method/path pair
    - returns none when no handler is registered
    # dispatch
