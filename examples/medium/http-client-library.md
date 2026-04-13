# Requirement: "an ergonomic HTTP client library"

A fluent HTTP client for building requests, setting headers and query params, sending, and reading typed response bodies.

std
  std.http
    std.http.send
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_raw_response, string]
      + performs the request and returns status, headers, and body
      - returns error on network failure
      # http
  std.strings
    std.strings.join
      @ (parts: list[string], sep: string) -> string
      + joins parts with sep
      # strings
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

httpclient
  httpclient.new_request
    @ (method: string, url: string) -> request_state
    + returns a request builder with no headers, no query, and an empty body
    # construction
  httpclient.set_header
    @ (req: request_state, name: string, value: string) -> request_state
    + sets a header, replacing any existing value for the same name
    # builder
  httpclient.add_query
    @ (req: request_state, key: string, value: string) -> request_state
    + appends a query parameter
    + permits multiple values for the same key
    # builder
  httpclient.set_body
    @ (req: request_state, body: bytes) -> request_state
    + sets the request body
    # builder
  httpclient.build_url
    @ (req: request_state) -> string
    + returns the final URL including query string
    + url-encodes key and value
    # builder
    -> std.strings.join
  httpclient.send
    @ (req: request_state) -> result[response, string]
    + returns a response with status, headers, and body
    - returns error on network failure
    # execution
    -> httpclient.build_url
    -> std.http.send
  httpclient.response_json
    @ (resp: response) -> result[map[string, string], string]
    + returns the body parsed as a JSON object
    - returns error when the body is not valid JSON
    # decoding
    -> std.json.parse_object
