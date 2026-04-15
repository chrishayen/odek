# Requirement: "a promise-based HTTP client that works in both server and browser environments"

Exposes an async HTTP client with request/response interceptors and automatic body serialization.

std
  std.io
    std.io.http_send
      fn (method: string, url: string, headers: map[string,string], body: bytes, timeout_ms: i32) -> result[http_raw_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error on network failure or timeout
      # http
  std.encoding
    std.encoding.json_encode
      fn (value: json_value) -> string
      + serializes a JSON value to its text form
      # serialization
    std.encoding.json_decode
      fn (text: string) -> result[json_value, string]
      + parses JSON text into a JSON value
      - returns error on malformed JSON
      # serialization
    std.encoding.url_encode_query
      fn (params: map[string,string]) -> string
      + returns a percent-encoded query string
      # url
  std.async
    std.async.resolve
      fn (value: http_response) -> promise[http_response]
      + creates a promise already resolved with the given value
      # async
    std.async.reject
      fn (err: string) -> promise[http_response]
      + creates a promise already rejected with the given error
      # async
    std.async.run
      fn (fn: callable) -> promise[http_response]
      + runs a function on the async runtime and returns its promise
      # async

http_client
  http_client.new
    fn (base_url: string, default_headers: map[string,string], timeout_ms: i32) -> client_state
    + creates a client with base URL, default headers, and timeout
    # construction
  http_client.add_request_interceptor
    fn (state: client_state, fn: callable) -> client_state
    + registers a function invoked on each outgoing request
    # interceptors
  http_client.add_response_interceptor
    fn (state: client_state, fn: callable) -> client_state
    + registers a function invoked on each incoming response
    # interceptors
  http_client.build_request
    fn (state: client_state, method: string, path: string, query: map[string,string], body: optional[json_value]) -> http_request
    + merges base URL, query string, headers, and serialized body
    # request_building
    -> std.encoding.url_encode_query
    -> std.encoding.json_encode
  http_client.request
    fn (state: client_state, req: http_request) -> promise[http_response]
    + runs request interceptors, sends the request, runs response interceptors
    + parses JSON bodies when the content-type header indicates JSON
    - rejects with an error message on non-2xx status
    - rejects on network failure or timeout
    # request_execution
    -> std.async.run
    -> std.async.resolve
    -> std.async.reject
    -> std.io.http_send
    -> std.encoding.json_decode
  http_client.get
    fn (state: client_state, path: string, query: map[string,string]) -> promise[http_response]
    + sends a GET request with the given query parameters
    # convenience
  http_client.post
    fn (state: client_state, path: string, body: json_value) -> promise[http_response]
    + sends a POST request with a JSON body
    # convenience
