# Requirement: "a rest api client library"

A fluent builder that accumulates base URL, path segments, and query parameters, then executes through an HTTP primitive.

std
  std.http
    std.http.get
      fn (url: string) -> result[string, string]
      + returns the response body for 2xx status codes
      - returns error on network failure or non-2xx status
      # http
    std.http.post
      fn (url: string, body: string) -> result[string, string]
      + returns the response body for 2xx status codes
      - returns error on network failure or non-2xx status
      # http
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

rest_client
  rest_client.new
    fn (base_url: string) -> rest_client_state
    + returns a client rooted at the base URL
    # construction
  rest_client.resource
    fn (client: rest_client_state, segment: string) -> rest_client_state
    + returns a client with the segment appended to the path
    # path
  rest_client.with_query
    fn (client: rest_client_state, key: string, value: string) -> rest_client_state
    + returns a client with an additional query parameter
    # query
  rest_client.get
    fn (client: rest_client_state) -> result[map[string, string], string]
    + fetches the resource and parses the JSON response
    - returns error on request failure or invalid JSON
    # request
    -> std.http.get
    -> std.json.parse_object
  rest_client.post
    fn (client: rest_client_state, body: string) -> result[map[string, string], string]
    + posts the body and parses the JSON response
    - returns error on request failure or invalid JSON
    # request
    -> std.http.post
    -> std.json.parse_object
