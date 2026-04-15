# Requirement: "an asynchronous HTTP client and server framework"

A full async HTTP stack: low-level socket I/O and HTTP parsing live in std; the project exposes an ergonomic client and a route-based server.

std
  std.net
    std.net.tcp_listen
      fn (host: string, port: u16) -> result[listener_state, string]
      + binds and listens on the given address
      - returns error when the port is already in use
      # networking
    std.net.tcp_accept
      fn (listener: listener_state) -> result[conn_state, string]
      + blocks on the event loop until a new connection arrives
      # networking
    std.net.tcp_dial
      fn (host: string, port: u16) -> result[conn_state, string]
      + opens a TCP connection to the target
      - returns error on DNS failure or connection refused
      # networking
    std.net.conn_read
      fn (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes, suspending the task if none are ready
      # networking
    std.net.conn_write
      fn (conn: conn_state, data: bytes) -> result[i32, string]
      + writes all bytes, suspending until the socket is drained
      # networking
    std.net.conn_close
      fn (conn: conn_state) -> void
      # networking
  std.async
    std.async.spawn
      fn (task: task_fn) -> task_handle
      + schedules a task on the event loop
      # concurrency
    std.async.run_loop
      fn () -> void
      + drives the event loop until all tasks complete
      # concurrency
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from a raw HTTP/1.1 message
      - returns error on malformed request line or headers
      # parsing
    std.http.parse_response
      fn (raw: bytes) -> result[http_response, string]
      + parses status, headers, and body from a raw HTTP/1.1 response
      - returns error when Content-Length is inconsistent with body
      # parsing
    std.http.encode_request
      fn (method: string, path: string, headers: map[string,string], body: bytes) -> bytes
      + serializes a request to HTTP/1.1 wire format
      # serialization
    std.http.encode_response
      fn (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes a response to HTTP/1.1 wire format
      # serialization
  std.url
    std.url.parse
      fn (raw: string) -> result[url_parts, string]
      + splits a URL into scheme, host, port, path, query
      - returns error on missing scheme
      # parsing

httpx
  httpx.client_new
    fn () -> client_state
    + creates a reusable client with a connection pool
    # construction
  httpx.client_request
    fn (client: client_state, method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
    + sends the request and awaits the full response
    - returns error on unreachable host or malformed URL
    # client_request
    -> std.url.parse
    -> std.net.tcp_dial
    -> std.http.encode_request
    -> std.net.conn_write
    -> std.net.conn_read
    -> std.http.parse_response
  httpx.server_new
    fn (host: string, port: u16) -> result[server_state, string]
    + prepares a server bound to host:port
    # construction
    -> std.net.tcp_listen
  httpx.server_route
    fn (server: server_state, method: string, path: string, handler: handler_fn) -> server_state
    + registers a handler for a (method, path) pair
    # routing
  httpx.server_serve
    fn (server: server_state) -> result[void, string]
    + accepts connections in a loop, parsing requests and dispatching to handlers
    - returns error when the listener is closed unexpectedly
    # serving
    -> std.net.tcp_accept
    -> std.async.spawn
    -> std.http.parse_request
    -> std.http.encode_response
    -> std.net.conn_write
    -> std.net.conn_close
  httpx.server_run
    fn (server: server_state) -> result[void, string]
    + drives the event loop until the server is stopped
    # lifecycle
    -> std.async.run_loop
