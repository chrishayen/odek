# Requirement: "a lightweight and fast HTTP router"

Routes HTTP requests to registered handlers by method and path pattern. Pattern matching supports static segments and named parameters.

std
  std.string
    std.string.split
      @ (s: string, sep: string) -> list[string]
      + splits on separator returning all segments
      + returns single-element list when separator is absent
      # strings

router
  router.new
    @ () -> router_state
    + returns an empty router
    # construction
  router.add
    @ (state: router_state, method: string, pattern: string, handler_id: string) -> router_state
    + registers a handler for the given method and pattern
    + pattern segments beginning with ':' become named parameters
    # registration
  router.find
    @ (state: router_state, method: string, path: string) -> optional[tuple[string, map[string, string]]]
    + returns (handler_id, params) when a route matches
    - returns none when no pattern matches the method/path pair
    # dispatch
    -> std.string.split
