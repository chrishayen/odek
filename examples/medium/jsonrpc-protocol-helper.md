# Requirement: "a JSON-RPC 2.0 protocol helper library"

Encode and decode request, response, and error envelopes. Callers still choose their own transport.

std
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on invalid JSON
      # serialization

jsonrpc
  jsonrpc.encode_request
    fn (id: string, method: string, params: map[string, string]) -> string
    + returns a JSON-RPC 2.0 request envelope
    # request_framing
    -> std.json.encode_object
  jsonrpc.encode_notification
    fn (method: string, params: map[string, string]) -> string
    + returns a JSON-RPC 2.0 notification (no id field)
    # request_framing
    -> std.json.encode_object
  jsonrpc.encode_response
    fn (id: string, result_value: map[string, string]) -> string
    + returns a JSON-RPC 2.0 success response
    # response_framing
    -> std.json.encode_object
  jsonrpc.encode_error
    fn (id: string, code: i32, message: string) -> string
    + returns a JSON-RPC 2.0 error response
    # response_framing
    -> std.json.encode_object
  jsonrpc.decode_message
    fn (raw: string) -> result[map[string, string], string]
    + parses an incoming envelope and returns its fields
    - returns error when the envelope is missing the jsonrpc field
    - returns error when the envelope is not a JSON object
    # parsing
    -> std.json.parse_object
  jsonrpc.classify
    fn (message: map[string, string]) -> string
    + returns one of "request", "notification", "response", "error"
    # parsing
