# Requirement: "a resource-oriented web framework"

Routes are derived from resource definitions: each resource declares its fields and which operations (list, get, create, update, delete) it supports.

std
  std.http
    std.http.parse_request
      @ (raw: string) -> result[http_request, string]
      + parses an HTTP request
      - returns error on malformed input
      # http
    std.http.render_response
      @ (status: i32, body: string) -> string
      + returns a wire-format response
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on malformed input
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

resource_framework
  resource_framework.define_resource
    @ (name: string, fields: list[string], operations: list[string]) -> resource_definition
    + returns a resource descriptor with the allowed operations
    # schema
  resource_framework.routes_for
    @ (resource: resource_definition) -> list[route_entry]
    + returns the conventional routes for the resource's operations
    + list maps to GET /{name}, create to POST /{name}, and so on
    # routing
  resource_framework.dispatch
    @ (resources: list[resource_definition], handlers: resource_handlers, request: http_request) -> http_response
    + routes the request to the matching operation handler
    - returns 404 when no resource matches
    - returns 405 when the operation is not supported
    # dispatch
  resource_framework.validate_payload
    @ (resource: resource_definition, payload: map[string, string]) -> result[void, string]
    + returns ok when every declared field is present
    - returns error listing missing fields
    # validation
  resource_framework.handle
    @ (resources: list[resource_definition], handlers: resource_handlers, raw_request: string) -> string
    + parses, dispatches, and renders a full request-response cycle
    - returns 400 when the request is malformed
    # dispatch
    -> std.http.parse_request
    -> std.http.render_response
    -> std.json.parse_object
    -> std.json.encode_object
