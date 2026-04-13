# Requirement: "a modular, observable framework for backend services"

Wires together configuration, dependency injection, http routing, health checks, logging, metrics, and tracing around a service lifecycle.

std
  std.log
    std.log.info
      @ (message: string, fields: map[string, string]) -> void
      + emits an informational structured log entry
      # logging
    std.log.error
      @ (message: string, fields: map[string, string]) -> void
      + emits an error structured log entry
      # logging
  std.metrics
    std.metrics.counter_inc
      @ (name: string, labels: map[string, string]) -> void
      + increments a named counter by 1
      # metrics
    std.metrics.histogram_observe
      @ (name: string, value: f64, labels: map[string, string]) -> void
      + records a value into a named histogram
      # metrics
  std.trace
    std.trace.start_span
      @ (name: string) -> span_handle
      + starts a new span and returns its handle
      # tracing
    std.trace.end_span
      @ (span: span_handle) -> void
      + ends the given span
      # tracing
  std.http
    std.http.new_router
      @ () -> http_router
      + creates an empty http router
      # http
    std.http.route
      @ (router: http_router, method: string, path: string, handler: http_handler) -> http_router
      + binds a handler to the given method and path
      # http
    std.http.listen
      @ (router: http_router, addr: string) -> result[http_server, string]
      + binds the router to addr and starts serving
      - returns error when addr is already in use
      # http
    std.http.shutdown
      @ (server: http_server) -> result[void, string]
      + drains in-flight requests and stops the server
      # http
  std.config
    std.config.load_env
      @ (prefix: string) -> map[string, string]
      + returns environment variables with the given prefix, with the prefix stripped
      # config

yokai
  yokai.new_app
    @ (name: string) -> app_state
    + creates an empty application with a service name
    # construction
  yokai.register
    @ (app: app_state, key: string, provider: component_provider) -> app_state
    + registers a component provider in the dependency graph
    # di
  yokai.resolve
    @ (app: app_state, key: string) -> result[component_handle, string]
    + instantiates and returns the component for key, memoizing the result
    - returns error when key is not registered
    - returns error when the provider returns an error
    # di
  yokai.load_config
    @ (app: app_state, env_prefix: string) -> app_state
    + loads configuration from environment variables with the given prefix
    # config
    -> std.config.load_env
  yokai.mount_route
    @ (app: app_state, method: string, path: string, handler: http_handler) -> app_state
    + registers an http route on the application's router
    # http
    -> std.http.route
  yokai.add_health_check
    @ (app: app_state, name: string, check: health_check) -> app_state
    + registers a named liveness/readiness probe
    # health
  yokai.run_health_checks
    @ (app: app_state) -> health_report
    + runs every registered probe and returns a combined report
    # health
    -> std.log.info
  yokai.with_logging
    @ (handler: http_handler) -> http_handler
    + wraps a handler to emit structured access logs
    # middleware
    -> std.log.info
    -> std.log.error
  yokai.with_metrics
    @ (handler: http_handler) -> http_handler
    + wraps a handler to count requests and time latency
    # middleware
    -> std.metrics.counter_inc
    -> std.metrics.histogram_observe
  yokai.with_tracing
    @ (handler: http_handler) -> http_handler
    + wraps a handler in a tracing span
    # middleware
    -> std.trace.start_span
    -> std.trace.end_span
  yokai.start
    @ (app: app_state, addr: string) -> result[app_state, string]
    + resolves every component and starts serving on addr
    - returns error when any provider fails
    # lifecycle
    -> std.http.new_router
    -> std.http.listen
  yokai.stop
    @ (app: app_state) -> result[void, string]
    + drains the http server and tears down components in reverse registration order
    # lifecycle
    -> std.http.shutdown
