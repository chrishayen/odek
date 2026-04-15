# Requirement: "a declarative end-to-end http api testing library"

Builds requests fluently and asserts over status, headers, and json bodies.

std
  std.http
    std.http.send
      fn (method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + returns the response object including status, headers, and body
      - returns error on transport failure
      # networking
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a json document into a generic value
      - returns error on invalid json
      # serialization
    std.json.path
      fn (root: json_value, dotted_path: string) -> optional[json_value]
      + returns the value at the given dotted path
      # serialization

http_expect
  http_expect.new
    fn (base_url: string) -> expect_state
    + creates a test client bound to a base url
    # construction
  http_expect.request
    fn (state: expect_state, method: string, path: string) -> request_builder
    + begins a request for the given method and path
    # request
  http_expect.with_header
    fn (b: request_builder, name: string, value: string) -> request_builder
    + attaches a header to the pending request
    # request
  http_expect.with_json_body
    fn (b: request_builder, body: string) -> request_builder
    + sets the request body to json and the content type header
    # request
  http_expect.execute
    fn (b: request_builder) -> result[assertion_state, string]
    + sends the request and returns an assertable response
    - returns error on transport failure
    # execution
    -> std.http.send
  http_expect.expect_status
    fn (a: assertion_state, code: i32) -> result[assertion_state, string]
    + passes when the response status equals code
    - returns error with a descriptive mismatch message
    # assertion
  http_expect.expect_header
    fn (a: assertion_state, name: string, value: string) -> result[assertion_state, string]
    + passes when the named response header equals value
    - returns error when header is missing or differs
    # assertion
  http_expect.expect_json_path
    fn (a: assertion_state, dotted_path: string, expected: string) -> result[assertion_state, string]
    + passes when the value at dotted_path in the json body equals expected
    - returns error when the path is missing
    # assertion
    -> std.json.parse
    -> std.json.path
