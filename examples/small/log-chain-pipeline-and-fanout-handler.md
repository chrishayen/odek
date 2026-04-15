# Requirement: "a library for chaining structured log handlers into pipelines and fanouts"

Handlers are composable: one can forward a record through a middleware chain or broadcast to multiple sinks.

std: (all units exist)

log_chain
  log_chain.pipeline
    fn (middleware: list[log_middleware], terminal: log_handler) -> log_handler
    + returns a handler that runs each middleware in order before invoking terminal
    ? middleware may transform or drop the record; dropping short-circuits the chain
    # composition
  log_chain.fanout
    fn (handlers: list[log_handler]) -> log_handler
    + returns a handler that forwards each record to every handler in the list
    + continues to remaining handlers when one returns an error
    # composition
  log_chain.filter
    fn (predicate: fn(log_record) -> bool) -> log_middleware
    + returns middleware that passes records matching the predicate and drops the rest
    # middleware
  log_chain.map
    fn (transform: fn(log_record) -> log_record) -> log_middleware
    + returns middleware that rewrites each record before forwarding
    # middleware
