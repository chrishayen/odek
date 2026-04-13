# Requirement: "a REST API framework that connects to pluggable backend data sources"

A framework that lets callers declare resources backed by arbitrary data sources and exposes CRUD-style HTTP endpoints. Framework builds routing and dispatch on top of HTTP and JSON primitives.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body
      - returns error on malformed input
      # http
    std.http.build_response
      @ (status: i32, headers: map[string, string], body: bytes) -> bytes
      + serializes a response into wire bytes
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on non-object input
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string map as a JSON object
      # serialization
    std.json.encode_array
      @ (items: list[map[string, string]]) -> string
      + encodes a list of string maps as a JSON array
      # serialization

apiframe
  apiframe.new_app
    @ () -> app_state
    + returns an empty application
    # construction
  apiframe.register_resource
    @ (app: app_state, name: string, source_id: string) -> app_state
    + maps a resource name to a data source identifier
    # registration
  apiframe.bind_source
    @ (app: app_state, source_id: string, source: data_source) -> app_state
    + attaches a concrete data source for the given id
    ? data_source is an opaque handle whose implementation is passed by the caller
    # registration
  apiframe.route
    @ (app: app_state, method: string, path: string) -> optional[resource_op]
    + resolves a method and path like "GET /users/42" to a resource + op + id
    - returns none when no resource matches
    # routing
  apiframe.handle_list
    @ (app: app_state, resource: string) -> result[list[map[string, string]], string]
    + fetches all records from the bound source
    - returns error when the resource is unbound
    # handler
  apiframe.handle_get
    @ (app: app_state, resource: string, id: string) -> result[map[string, string], string]
    + fetches a single record by id
    - returns error when no record exists
    # handler
  apiframe.handle_create
    @ (app: app_state, resource: string, body: map[string, string]) -> result[string, string]
    + inserts a record and returns its new id
    # handler
  apiframe.dispatch
    @ (app: app_state, raw: bytes) -> bytes
    + parses the request, routes to the handler, and returns a response
    - returns a 404 response when no route matches
    - returns a 500 response when the handler errors
    # dispatch
    -> std.http.parse_request
    -> std.http.build_response
    -> std.json.parse_object
    -> std.json.encode_object
    -> std.json.encode_array
