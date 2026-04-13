# Requirement: "a small HTTP client with a fluent request builder that issues requests and parses responses"

Build a request value, send it, get a response with helpers for common body types.

std
  std.http
    std.http.send
      @ (method: string, url: string, headers: map[string,string], body: bytes, timeout_ms: i32) -> result[http_response, string]
      + performs an HTTP(S) request with the given timeout
      - returns error on network or timeout
      # http
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      - returns error on invalid JSON
      # parsing
    std.json.encode
      @ (value: json_value) -> string
      # serialization
  std.encoding
    std.encoding.url_encode_query
      @ (params: map[string,string]) -> string
      + encodes a map as a URL query string
      # encoding

http_client
  http_client.new
    @ (base_url: string) -> client_state
    + returns a client with a base URL applied to relative paths
    # construction
  http_client.with_header
    @ (state: client_state, name: string, value: string) -> client_state
    + adds a default header to every request
    # configuration
  http_client.with_timeout
    @ (state: client_state, ms: i32) -> client_state
    + sets the default request timeout in milliseconds
    # configuration
  http_client.build_request
    @ (state: client_state, method: string, path: string, query: map[string,string], body: bytes) -> prepared_request
    + resolves the URL and merges default headers
    # builder
    -> std.encoding.url_encode_query
  http_client.send
    @ (state: client_state, req: prepared_request) -> result[http_response, string]
    + sends the prepared request
    - returns error on network or timeout
    # send
    -> std.http.send
  http_client.get_json
    @ (state: client_state, path: string, query: map[string,string]) -> result[json_value, string]
    + convenience for GET + JSON parsing
    - returns error on non-2xx or parse failure
    # helpers
    -> http_client.build_request
    -> http_client.send
    -> std.json.parse
  http_client.post_json
    @ (state: client_state, path: string, body: json_value) -> result[json_value, string]
    + convenience for POST + JSON encode + JSON parse
    - returns error on non-2xx or parse failure
    # helpers
    -> http_client.build_request
    -> http_client.send
    -> std.json.encode
    -> std.json.parse
