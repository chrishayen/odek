# Requirement: "a readable, flexible client for REST APIs"

A fluent builder for HTTP requests that serializes JSON bodies and decodes JSON responses.

std
  std.http
    std.http.send
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs the HTTP request and returns status, headers, and body
      - returns error on network failure
      # http
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + serializes a JSON value to a string
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses JSON text into a value tree
      - returns error on invalid JSON
      # serialization

rest_client
  rest_client.new
    fn (base_url: string) -> rest_client_state
    + creates a client rooted at base_url
    # construction
  rest_client.with_header
    fn (client: rest_client_state, name: string, value: string) -> rest_client_state
    + returns a client with the header appended to every request
    # configuration
  rest_client.request
    fn (client: rest_client_state, method: string, path: string, body: optional[json_value]) -> result[json_value, string]
    + sends a request to base_url + path and decodes the JSON response
    + includes client-level headers; serializes body as JSON when present
    - returns error when the status code is not in the 2xx range
    - returns error when the response body is not valid JSON
    # request_execution
    -> std.json.encode
    -> std.http.send
    -> std.json.parse
