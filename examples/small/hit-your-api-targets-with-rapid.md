# Requirement: "a fluent HTTP client"

A builder-style HTTP client that composes a request and executes it through a std primitive.

std
  std.net
    std.net.http_request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[tuple[i32, map[string, string], bytes], string]
      + sends an HTTP request and returns status, headers, and body
      - returns error on connection failure or timeout
      # networking

http_client
  http_client.new
    @ (base_url: string) -> client_state
    + creates a client rooted at a base URL
    # construction
  http_client.header
    @ (c: client_state, name: string, value: string) -> client_state
    + adds a default header applied to every request
    # configuration
  http_client.build
    @ (c: client_state, method: string, path: string) -> request_builder
    + starts a fluent request builder for a method and path
    # building
  http_client.with_query
    @ (rb: request_builder, key: string, value: string) -> request_builder
    + appends a query parameter
    # building
  http_client.with_body
    @ (rb: request_builder, body: bytes) -> request_builder
    + sets the request body
    # building
  http_client.send
    @ (rb: request_builder) -> result[tuple[i32, map[string, string], bytes], string]
    + executes the built request and returns the response
    - returns error on network failure
    # execution
    -> std.net.http_request
