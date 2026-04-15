# Requirement: "an HTTP client builder for composing API requests"

A fluent builder for HTTP requests plus a thin execute step. The builder is pure; network I/O lives in one std primitive.

std
  std.http
    std.http.send
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs a blocking request and returns status, headers, and body
      - returns error on transport failure
      # http
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

sling
  sling.new
    fn (base_url: string) -> sling_state
    + creates a builder with the given base URL and an empty method, headers, and query
    # construction
  sling.method
    fn (state: sling_state, method: string) -> sling_state
    + sets the HTTP method on the builder
    # building
  sling.path
    fn (state: sling_state, path: string) -> sling_state
    + appends path to the base URL, collapsing adjacent slashes
    # building
  sling.set_header
    fn (state: sling_state, name: string, value: string) -> sling_state
    + sets or replaces a header on the builder
    # building
  sling.query_param
    fn (state: sling_state, name: string, value: string) -> sling_state
    + appends a query parameter to the builder
    # building
  sling.json_body
    fn (state: sling_state, obj: map[string, string]) -> sling_state
    + sets the body to the JSON encoding of obj and the Content-Type header to application/json
    # building
    -> std.json.encode_object
  sling.execute
    fn (state: sling_state) -> result[http_response, string]
    + assembles the final URL and sends the request
    - returns error when method is unset
    # dispatch
    -> std.http.send
