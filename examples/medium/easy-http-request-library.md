# Requirement: "an ergonomic HTTP request library"

A small fluent layer over a low-level request primitive with convenience helpers for common verbs and body types.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + issues an HTTP request and returns status and body
      - returns error on network failure
      # http
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
  std.url
    std.url.encode_query
      fn (params: map[string, string]) -> string
      + returns a URL-encoded query string
      # url

easy_http
  easy_http.get
    fn (url: string, params: map[string, string]) -> result[http_response, string]
    + appends params as a query string and issues a GET request
    # http
    -> std.http.request
    -> std.url.encode_query
  easy_http.post_json
    fn (url: string, body: map[string, string]) -> result[http_response, string]
    + sends a JSON body with the correct content type
    # http
    -> std.http.request
    -> std.json.encode_object
  easy_http.post_form
    fn (url: string, fields: map[string, string]) -> result[http_response, string]
    + sends a form-encoded body with the correct content type
    # http
    -> std.http.request
    -> std.url.encode_query
  easy_http.response_json
    fn (response: http_response) -> result[map[string, string], string]
    + parses a JSON object from the response body
    - returns error when the status code is not in the 2xx range
    # decoding
    -> std.json.parse_object
