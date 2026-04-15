# Requirement: "a library for chaining request handlers with per-request scoped data"

std: (all units exist)

handler_chain
  handler_chain.new
    fn () -> chain_state
    + creates an empty handler chain
    # construction
  handler_chain.use
    fn (state: chain_state, middleware: middleware_fn) -> chain_state
    + appends middleware to the chain; later middleware runs inside earlier middleware
    # composition
  handler_chain.finalize
    fn (state: chain_state, terminal: handler_fn) -> handler_fn
    + returns a single handler that invokes the chain around the terminal handler
    # composition
  handler_chain.scope_get
    fn (ctx: request_context, key: string) -> optional[string]
    + returns a value previously stored on the request context
    - returns none when the key is not present
    # scoped_state
  handler_chain.scope_set
    fn (ctx: request_context, key: string, value: string) -> request_context
    + returns a new context with the key set to value
    # scoped_state
