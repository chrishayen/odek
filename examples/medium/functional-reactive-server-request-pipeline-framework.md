# Requirement: "a functional-reactive server request pipeline framework"

Requests flow through a middleware pipeline of pure functions. Each step returns either the next context or a short-circuit response.

std: (all units exist)

pipeline
  pipeline.new_context
    fn (method: string, path: string, headers: map[string, string], body: bytes) -> request_context
    + creates an immutable request context
    # construction
  pipeline.compose
    fn (middlewares: list[middleware_fn]) -> middleware_fn
    + returns a single middleware that runs each input middleware in order
    + a short-circuit response from any stage stops downstream stages
    # composition
  pipeline.map_response
    fn (ctx: request_context, f: func(response) -> response) -> middleware_fn
    + returns a middleware that rewrites the downstream response via f
    # transform
  pipeline.match_path
    fn (pattern: string, handler: middleware_fn) -> middleware_fn
    + returns a middleware that runs handler when the pattern matches ctx.path
    + pattern supports ":param" placeholders and writes them to ctx
    # routing
  pipeline.match_method
    fn (method: string, handler: middleware_fn) -> middleware_fn
    + returns a middleware that runs handler only when ctx.method matches
    - other methods pass through without running the handler
    # routing
  pipeline.run
    fn (mw: middleware_fn, ctx: request_context) -> response
    + executes the middleware and returns the produced response
    + returns a 404 response when no middleware produces one
    # execution
