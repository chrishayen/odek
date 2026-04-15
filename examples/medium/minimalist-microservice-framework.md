# Requirement: "a minimalist microservice framework"

Services are collections of named handlers that take a request payload and return a response. The framework handles routing and error envelopes. Transport is the caller's responsibility.

std
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization

service
  service.new
    fn (name: string) -> service_state
    + returns an empty service with the given name
    # construction
  service.register_handler
    fn (state: service_state, route: string, handler_id: string) -> service_state
    + binds an opaque handler id to a route
    ? the framework is handler-agnostic; dispatch returns the id
    # registration
  service.dispatch
    fn (state: service_state, route: string, raw_body: string) -> result[string, string]
    + parses the body, looks up the route, and returns the handler id as a signal for the caller to invoke
    - returns error envelope when the route is unknown
    - returns error envelope when the body is not valid JSON
    # routing
    -> std.json.parse_object
  service.ok_response
    fn (data: map[string, string]) -> string
    + encodes a success envelope: {"ok": "true", ...data}
    # responses
    -> std.json.encode_object
  service.error_response
    fn (message: string) -> string
    + encodes an error envelope: {"ok": "false", "error": message}
    # responses
    -> std.json.encode_object
