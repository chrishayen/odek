# Requirement: "a micro-framework for writing RESTful API integration tests"

Composable assertion helpers over an HTTP response so tests read like fluent expectations.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error when the connection fails
      # networking
  std.json
    std.json.get_path
      fn (raw: string, path: string) -> result[string, string]
      + returns the string value at a dotted JSON path
      - returns error when the path does not exist
      # serialization

restit
  restit.call
    fn (method: string, url: string, body: bytes) -> result[test_response, string]
    + issues a request and wraps the response for assertions
    # request
    -> std.http.request
  restit.expect_status
    fn (resp: test_response, code: i32) -> result[test_response, string]
    + passes through when the status matches
    - returns an error describing the mismatch otherwise
    # assertion
  restit.expect_header
    fn (resp: test_response, name: string, value: string) -> result[test_response, string]
    + passes through when the header equals the expected value
    - returns an error with the observed value otherwise
    # assertion
  restit.expect_json_path
    fn (resp: test_response, path: string, value: string) -> result[test_response, string]
    + passes through when the JSON path matches the expected value
    - returns an error when the path is missing or differs
    # assertion
    -> std.json.get_path
