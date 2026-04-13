# Requirement: "an HTTP client with a plugin pipeline for middleware"

Requests pass through a chain of registered middleware that can mutate them before sending and mutate responses on the way back. The wire send is a std primitive.

std
  std.http
    std.http.send
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs the request and returns status, headers, and body
      - returns error on transport failure
      # http

http_client
  http_client.new_client
    @ () -> client_state
    + returns a client with an empty middleware chain
    # construction
  http_client.use_header
    @ (client: client_state, name: string, value: string) -> client_state
    + appends a middleware that adds the header to every outgoing request
    # middleware
  http_client.use_retry
    @ (client: client_state, max_attempts: i32) -> client_state
    + appends a middleware that retries on transport errors up to max_attempts
    # middleware
  http_client.use_base_url
    @ (client: client_state, base: string) -> client_state
    + appends a middleware that prepends base to relative request urls
    # middleware
  http_client.new_request
    @ (method: string, url: string) -> request_state
    + returns a request with the given method and url and empty headers and body
    # construction
  http_client.set_body
    @ (request: request_state, body: bytes) -> request_state
    + sets the request body
    # construction
  http_client.send
    @ (client: client_state, request: request_state) -> result[http_response, string]
    + runs the middleware chain and returns the final response
    - returns error when transport fails after all retries
    # execution
    -> std.http.send
