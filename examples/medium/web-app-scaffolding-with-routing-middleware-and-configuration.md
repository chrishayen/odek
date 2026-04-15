# Requirement: "a web application scaffolding library with routing, middleware, and configuration"

Thin layered pieces a web service wires together: router, middleware chain, config loader, request context.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the file does not exist
      # io
  std.env
    std.env.get
      fn (key: string) -> optional[string]
      + returns the environment variable when set
      # environment

web_app
  web_app.new_app
    fn (name: string) -> app_state
    + returns an empty app with the given name
    # construction
  web_app.load_config
    fn (path: string) -> result[config_state, string]
    + reads a key=value config file and overlays environment variables
    - returns error when the file is malformed
    # configuration
    -> std.fs.read_all
    -> std.env.get
  web_app.route
    fn (a: app_state, method: string, path: string, handler: handler_fn) -> app_state
    + registers a handler for an exact method+path pair
    - rejects registration when the pair already exists
    # routing
  web_app.route_param
    fn (a: app_state, method: string, pattern: string, handler: handler_fn) -> app_state
    + registers a handler for a pattern with ":name" segments
    # routing
  web_app.use
    fn (a: app_state, mw: middleware_fn) -> app_state
    + appends middleware to the chain; runs in registration order
    # middleware
  web_app.dispatch
    fn (a: app_state, req: request_data) -> response_data
    + runs middleware then the matching handler and returns its response
    - returns 404 when no route matches
    # execution
  web_app.context_value
    fn (req: request_data, key: string) -> optional[string]
    + returns a value placed on the request by earlier middleware
    # context
  web_app.with_context
    fn (req: request_data, key: string, value: string) -> request_data
    + returns a copy of the request with the key/value set
    # context
