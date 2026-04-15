# Requirement: "a high-performance http client"

An HTTP/1.1 client with connection pooling, request building, and response parsing. The transport layer is pluggable so the library can be tested without real sockets.

std
  std.net
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[conn_state, string]
      + opens a TCP connection
      - returns error when the host is unreachable
      # networking
    std.net.send_bytes
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes bytes to an open connection
      - returns error on a closed or broken connection
      # networking
    std.net.recv_until
      fn (conn: conn_state, delimiter: bytes, max_bytes: i32) -> result[bytes, string]
      + reads up to max_bytes or until the delimiter appears
      - returns error on a broken connection
      # networking
  std.text
    std.text.split
      fn (s: string, sep: string) -> list[string]
      + splits a string on a separator
      # text

httpc
  httpc.client_new
    fn (max_conns_per_host: i32) -> client_state
    + returns a client with an empty connection pool
    # construction
  httpc.request_build
    fn (method: string, url: string) -> result[request_state, string]
    + parses the URL and returns a request builder
    - returns error on a malformed URL
    # request_building
    -> std.text.split
  httpc.request_set_header
    fn (req: request_state, name: string, value: string) -> request_state
    + sets or replaces a header
    # request_building
  httpc.request_set_body
    fn (req: request_state, body: bytes) -> request_state
    + sets the body and updates Content-Length
    # request_building
  httpc.send
    fn (client: client_state, req: request_state) -> result[tuple[response, client_state], string]
    + sends the request, reusing a pooled connection when possible
    - returns error on network or parse failure
    # execution
    -> std.net.dial_tcp
    -> std.net.send_bytes
    -> std.net.recv_until
  httpc.parse_response
    fn (raw: bytes) -> result[response, string]
    + parses status line, headers, and body
    - returns error on malformed HTTP
    # response_parsing
    -> std.text.split
  httpc.release_connection
    fn (client: client_state, conn: conn_state) -> client_state
    + returns a live connection to the pool for reuse
    # pooling
  httpc.close
    fn (client: client_state) -> void
    + closes every pooled connection
    # lifecycle
