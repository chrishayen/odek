# Requirement: "a user-friendly HTTP client library"

A small wrapper around HTTP primitives that accepts headers and body as simple maps and returns a structured response.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + performs the request and returns status, headers, and body
      - returns error when the URL cannot be resolved
      # http
  std.json
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + serializes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

http_client
  http_client.get
    fn (url: string, headers: map[string,string]) -> result[http_response, string]
    + issues a GET and returns the response
    - returns error on transport failure
    # get
    -> std.http.request
  http_client.post_json
    fn (url: string, headers: map[string,string], body: map[string,string]) -> result[http_response, string]
    + serializes body as JSON, sets Content-Type, and issues a POST
    # post_json
    -> std.json.encode_object
    -> std.http.request
  http_client.decode_json
    fn (resp: http_response) -> result[map[string,string], string]
    + parses the response body as a JSON object
    - returns error when the status is not in the 2xx range
    # decoding
    -> std.json.parse_object
