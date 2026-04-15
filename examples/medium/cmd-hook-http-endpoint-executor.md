# Requirement: "a library for defining http endpoints that execute commands on the server"

Users register hooks mapping URL paths to command templates. On request, the library matches the path, substitutes request data into the template, runs the command, and returns its output as the response body.

std
  std.process
    std.process.run
      fn (command: string, args: list[string]) -> result[process_result, string]
      + runs the command and returns exit code, stdout, and stderr
      - returns error when the command cannot be launched
      # process
  std.http
    std.http.serve
      fn (port: u16, handler: fn(http_request) -> http_response) -> result[void, string]
      + starts an HTTP server that invokes the handler for each request
      - returns error when the port cannot be bound
      # http
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

cmd_hook
  cmd_hook.new
    fn () -> hook_registry
    + creates an empty hook registry
    # construction
  cmd_hook.register
    fn (reg: hook_registry, path: string, command: string, arg_templates: list[string]) -> hook_registry
    + adds a hook binding the path to a command whose argument templates may reference request fields
    # registration
  cmd_hook.render_args
    fn (templates: list[string], request: http_request) -> result[list[string], string]
    + returns one concrete argument per template by substituting header, query, and parsed body values
    - returns error when a template references a missing field
    # templating
    -> std.json.parse_object
  cmd_hook.dispatch
    fn (reg: hook_registry, request: http_request) -> http_response
    + returns a 200 response with the command's stdout when a matching hook runs successfully
    - returns 404 when no hook matches the request path
    - returns 500 when the command exits non-zero or template rendering fails
    # dispatch
    -> std.process.run
  cmd_hook.serve
    fn (reg: hook_registry, port: u16) -> result[void, string]
    + starts an HTTP server that dispatches every request through the registry
    - returns error when the port cannot be bound
    # server
    -> std.http.serve
