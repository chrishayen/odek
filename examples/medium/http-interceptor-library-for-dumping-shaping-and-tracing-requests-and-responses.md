# Requirement: "a library for adding interceptors to an HTTP client for dumping, shaping, and tracing requests and responses"

A middleware chain wrapping an HTTP round-tripper. Each interceptor sees the request, may mutate it, then sees the response.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.io
    std.io.write_string
      fn (sink: string, data: string) -> result[void, string]
      + appends data to a named sink (file path or logical stream)
      - returns error when the sink cannot be opened
      # io

http_interceptors
  http_interceptors.new_chain
    fn () -> chain_state
    + returns an empty interceptor chain
    # construction
  http_interceptors.add
    fn (chain: chain_state, name: string, before: string, after: string) -> chain_state
    + appends an interceptor identified by name with hook tags for before/after
    ? before/after are tag strings the execute step dispatches on
    # registration
  http_interceptors.execute
    fn (chain: chain_state, request: http_request) -> result[http_response, string]
    + runs before-hooks in order, performs the request, then runs after-hooks in reverse
    - returns error when any hook returns an error
    # dispatch
    -> std.time.now_millis
  http_interceptors.dump_hook
    fn (sink: string) -> interceptor_fn
    + returns an interceptor that writes the serialized request and response to a sink
    # builtin_hooks
    -> std.io.write_string
  http_interceptors.trace_hook
    fn () -> interceptor_fn
    + returns an interceptor that records timestamps on entry and exit
    # builtin_hooks
    -> std.time.now_millis
  http_interceptors.shape_hook
    fn (max_bytes: i64) -> interceptor_fn
    + returns an interceptor that truncates request bodies exceeding max_bytes
    - returns error when max_bytes is negative
    # builtin_hooks
