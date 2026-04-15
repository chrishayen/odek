# Requirement: "an HTTP middleware chain"

Composes a sequence of middleware handlers around a terminal handler, producing a single handler that runs them in order.

std: (all units exist)

middleware
  middleware.new_chain
    fn () -> chain_state
    + returns an empty chain
    # construction
  middleware.use
    fn (chain: chain_state, middleware_id: string) -> chain_state
    + appends a middleware to the chain
    # composition
  middleware.then
    fn (chain: chain_state, handler_id: string) -> handler_id
    + returns a handler that runs every middleware in order before the terminal handler
    # composition
  middleware.run
    fn (handler: handler_id, request: http_request) -> http_response
    + executes the composed chain for a request and returns the response
    # dispatch
