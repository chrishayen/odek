# Requirement: "a toolkit for building REST-style web APIs with serializers, routing, and authentication"

The project layer provides routing, request parsing, serialization, and auth middleware. std supplies the HTTP, JSON, and crypto primitives.

std
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from a raw request
      - returns error on malformed requests
      # http
    std.http.build_response
      fn (status: i32, headers: map[string, string], body: bytes) -> bytes
      + serializes a response into wire bytes
      # http
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON
      - returns error on invalid JSON
      # serialization
    std.json.encode_value
      fn (v: json_value) -> string
      + encodes a json value to a string
      # serialization
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + returns the 32-byte HMAC-SHA256 tag
      # cryptography

webapi
  webapi.new_router
    fn () -> router_state
    + returns an empty router
    # construction
  webapi.add_route
    fn (r: router_state, method: string, pattern: string, handler_id: string) -> router_state
    + registers a handler id for a method and path pattern with {param} segments
    # routing
  webapi.match_route
    fn (r: router_state, method: string, path: string) -> optional[route_match]
    + returns the handler id and extracted path parameters when a route matches
    - returns none when no route matches
    # routing
  webapi.parse_request
    fn (raw: bytes) -> result[http_request, string]
    + parses an incoming HTTP request
    # ingress
    -> std.http.parse_request
  webapi.decode_json_body
    fn (req: http_request) -> result[json_value, string]
    + parses the request body as JSON
    - returns error on invalid JSON
    # ingress
    -> std.json.parse_value
  webapi.serialize_resource
    fn (fields: map[string, json_value], visible: list[string]) -> string
    + produces a JSON object containing only the visible fields
    ? used by view-layer serializers to hide private fields
    # serialization
    -> std.json.encode_value
  webapi.json_response
    fn (status: i32, body: string) -> bytes
    + builds a JSON HTTP response with the given status
    # egress
    -> std.http.build_response
  webapi.error_response
    fn (status: i32, message: string) -> bytes
    + builds a JSON error response with a "error" field
    # egress
    -> std.http.build_response
  webapi.issue_token
    fn (user_id: string, secret: bytes) -> string
    + returns a signed token binding user_id to an HMAC tag
    # auth
    -> std.crypto.hmac_sha256
  webapi.verify_token
    fn (token: string, secret: bytes) -> result[string, string]
    + returns the user_id when the signature is valid
    - returns error when the token is malformed
    - returns error when the signature does not match
    # auth
    -> std.crypto.hmac_sha256
  webapi.require_auth
    fn (req: http_request, secret: bytes) -> result[string, string]
    + extracts the bearer token from the Authorization header and verifies it
    - returns error when the header is missing or invalid
    # auth
