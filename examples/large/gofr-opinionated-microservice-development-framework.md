# Requirement: "an opinionated microservice development framework"

A batteries-included service container: HTTP routing, configuration, structured logging, health checks, a data-source registry, and middleware chaining.

std
  std.http
    std.http.listen
      @ (port: i32, handler: fn(http_request) -> http_response) -> result[void, string]
      + binds a port and dispatches each incoming request to handler
      - returns error when the port is already bound
      # http_server
  std.env
    std.env.get
      @ (name: string) -> optional[string]
      + returns the value of an environment variable
      # config
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads an entire file as text
      - returns error when the file is missing
      # filesystem
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.json
    std.json.encode_object
      @ (fields: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization

service
  service.new_app
    @ (name: string) -> app_state
    + creates an app with default middleware and an empty router
    # construction
  service.load_config
    @ (app: app_state, config_path: string) -> result[app_state, string]
    + merges config from a file and environment variables, with env taking precedence
    - returns error when the config file is malformed
    # config
    -> std.fs.read_all
    -> std.env.get
  service.config_get
    @ (app: app_state, key: string) -> optional[string]
    + looks up a configuration value
    # config
  service.register_datasource
    @ (app: app_state, name: string, connect: fn() -> result[i32, string]) -> app_state
    + adds a named data source whose handle is opened lazily
    # datasource
  service.get_datasource
    @ (app: app_state, name: string) -> result[i32, string]
    + returns the cached handle for a registered data source, opening on first use
    - returns error when the data source is not registered
    # datasource
  service.use_middleware
    @ (app: app_state, mw: fn(http_request, fn(http_request) -> http_response) -> http_response) -> app_state
    + installs a middleware that wraps every handler
    # middleware
  service.route
    @ (app: app_state, method: string, path: string, handler: fn(http_request) -> http_response) -> app_state
    + registers a handler for a (method, path) pair
    # routing
  service.log_info
    @ (app: app_state, message: string, fields: map[string, string]) -> void
    + emits a structured info-level log record
    # logging
    -> std.time.now_millis
    -> std.json.encode_object
  service.log_error
    @ (app: app_state, message: string, fields: map[string, string]) -> void
    + emits a structured error-level log record
    # logging
    -> std.time.now_millis
    -> std.json.encode_object
  service.health_check
    @ (app: app_state) -> map[string, string]
    + returns a map of component name to "ok" or an error string, probing every registered data source
    # observability
  service.start
    @ (app: app_state, port: i32) -> result[void, string]
    + applies middleware, mounts the router, and begins accepting HTTP requests
    - returns error when no routes are registered
    # runtime
    -> std.http.listen
