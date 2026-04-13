# Requirement: "a proxy that converts JSON request bodies to protocol buffers"

Intercepts incoming requests, transcodes the body, and forwards the request upstream.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + returns a parsed json tree
      - returns error on malformed json
      # serialization
  std.http
    std.http.forward_request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + returns the upstream response
      - returns error when the connection fails
      # http

transcoder
  transcoder.load_schema
    @ (descriptor_bytes: bytes) -> result[schema_state, string]
    + returns a parsed protobuf descriptor usable for transcoding
    - returns error when the descriptor is malformed
    # schema
  transcoder.json_to_proto
    @ (schema: schema_state, message_name: string, raw_json: string) -> result[bytes, string]
    + returns the proto-encoded message for the given json
    - returns error when the message name is not in the schema
    - returns error when a required field is missing
    - returns error when a field type does not match
    # transcoding
    -> std.json.parse
  transcoder.route
    @ (path: string) -> result[string, string]
    + returns the message name configured for the given url path
    - returns error when the path has no mapping
    # routing
  transcoder.handle_request
    @ (schema: schema_state, upstream_base: string, method: string, path: string, headers: map[string, string], body: string) -> result[http_response, string]
    + returns the upstream response after transcoding the body
    - returns error when the path has no mapping
    - returns error when transcoding fails
    # proxy
    -> std.http.forward_request
