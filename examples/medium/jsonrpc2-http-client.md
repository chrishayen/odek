# Requirement: "a JSON-RPC 2.0 HTTP client"

The project exposes a single call entry point. Request framing, id generation, and response parsing are small dedicated helpers, with HTTP and JSON primitives in std.

std
  std.http
    std.http.post
      @ (url: string, headers: map[string, string], body: bytes) -> result[bytes, string]
      + sends a POST request and returns the response body
      - returns error on network failure or non-2xx status
      # http
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as a JSON object
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.random
    std.random.u64
      @ () -> u64
      + returns a random 64-bit unsigned integer
      # randomness

jsonrpc_client
  jsonrpc_client.new
    @ (endpoint: string) -> jsonrpc_client_state
    + creates a client bound to the given HTTP endpoint
    # construction
  jsonrpc_client.next_id
    @ (state: jsonrpc_client_state) -> tuple[string, jsonrpc_client_state]
    + returns a fresh request id and the advanced state
    # id_generation
    -> std.random.u64
  jsonrpc_client.build_request
    @ (id: string, method: string, params: map[string, string]) -> string
    + returns the JSON-RPC 2.0 request envelope as a string
    ? envelope fields are jsonrpc, method, params, id
    # request_framing
    -> std.json.encode_object
  jsonrpc_client.parse_response
    @ (raw: string) -> result[map[string, string], string]
    + returns the result object on success
    - returns error when the response carries a JSON-RPC error object
    - returns error when the envelope is malformed
    # response_parsing
    -> std.json.parse_object
  jsonrpc_client.call
    @ (state: jsonrpc_client_state, method: string, params: map[string, string]) -> result[map[string, string], string]
    + sends a request and returns the parsed result
    - returns error on transport or protocol failure
    # rpc_call
    -> std.http.post
