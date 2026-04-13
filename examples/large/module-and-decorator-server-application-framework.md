# Requirement: "a module-and-decorator-style server-side application framework"

A framework where callers register modules containing providers and controllers, resolve dependencies, and route HTTP requests to controller methods.

std
  std.net
    std.net.http_listen
      @ (host: string, port: u16, handler: fn(http_request) -> http_response) -> result[void, string]
      + starts an HTTP server that dispatches each request to handler
      - returns error when the port is already bound
      # http
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged value tree
      - returns error on malformed JSON
      # serialization
    std.json.encode
      @ (v: json_value) -> string
      + encodes a value tree as a JSON document
      # serialization

app_framework
  app_framework.new_module
    @ (name: string) -> module_def
    + creates an empty module with no providers, controllers, or imports
    # module
  app_framework.add_provider
    @ (m: module_def, token: string, factory: fn(list[any]) -> any, deps: list[string]) -> module_def
    + registers a provider keyed by token with its factory and its dependency tokens
    # dependency_injection
  app_framework.import_module
    @ (m: module_def, other: module_def) -> module_def
    + imports another module so its exported providers are visible
    # module
  app_framework.resolve
    @ (m: module_def) -> result[container, string]
    + topologically resolves all providers, instantiating each exactly once
    - returns error on a cyclic dependency
    - returns error when a provider references an unknown dependency token
    # dependency_injection
  app_framework.get
    @ (c: container, token: string) -> result[any, string]
    + returns the resolved instance for the given token
    - returns error when the token is not registered
    # dependency_injection
  app_framework.register_controller
    @ (m: module_def, path_prefix: string, routes: list[route_def]) -> module_def
    + registers a controller with a URL prefix and its route definitions
    # routing
  app_framework.build_router
    @ (m: module_def, c: container) -> router
    + produces a router that matches request method and path to a controller method
    # routing
  app_framework.handle
    @ (r: router, req: http_request) -> http_response
    + matches and invokes the appropriate controller method, returning its response
    + returns a 404 response when no route matches
    # routing
    -> std.json.parse
    -> std.json.encode
  app_framework.serve
    @ (r: router, host: string, port: u16) -> result[void, string]
    + binds an HTTP listener that feeds incoming requests to handle
    - returns error when binding fails
    # http
    -> std.net.http_listen
