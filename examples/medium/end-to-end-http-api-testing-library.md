# Requirement: "an end-to-end HTTP API testing library"

Builds a request, issues it, and runs a list of assertions against the response. Assertions return descriptive failures rather than booleans.

std
  std.http
    std.http.send_request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error on network failure
      # http
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document
      - returns error on invalid JSON
      # serialization
    std.json.pointer_lookup
      fn (value: json_value, pointer: string) -> optional[json_value]
      + returns the subvalue at a JSON pointer expression
      # serialization

apitest
  apitest.new_request
    fn (method: string, url: string) -> request_spec
    + creates a request with empty headers and body
    # construction
  apitest.with_header
    fn (req: request_spec, key: string, value: string) -> request_spec
    + sets a header on the request
    # construction
  apitest.with_json_body
    fn (req: request_spec, body: string) -> request_spec
    + sets the body to the given JSON string and content-type to application/json
    # construction
  apitest.send
    fn (req: request_spec) -> result[http_response, string]
    + executes the request
    # execution
    -> std.http.send_request
  apitest.expect_status
    fn (resp: http_response, expected: i32) -> optional[string]
    + returns none when status equals expected
    - returns a descriptive failure message when status differs
    # assertion
  apitest.expect_header
    fn (resp: http_response, name: string, expected: string) -> optional[string]
    + returns none when the named header matches expected
    - returns a descriptive failure when the header is missing or different
    # assertion
  apitest.expect_json_field
    fn (resp: http_response, pointer: string, expected: string) -> optional[string]
    + returns none when the JSON pointer resolves to expected
    - returns a descriptive failure when the pointer is missing or mismatched
    # assertion
    -> std.json.parse_value
    -> std.json.pointer_lookup
  apitest.run_assertions
    fn (resp: http_response, checks: list[assertion]) -> list[string]
    + returns all failure messages from running the given assertions
    + returns an empty list when every assertion passes
    # reporting
