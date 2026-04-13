# Requirement: "an HTTP server with routing"

Project layer is the user-facing API. Std carries all the real subsystems: HTTP parsing/formatting and TCP sockets — both reusable by any HTTP-using project.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses an HTTP/1.1 request from wire bytes
      + extracts method, path, headers, and body
      - returns error on a malformed request line
      - returns error on an incomplete request
      # http_parsing
    std.http.format_response
      @ (r: http_response) -> bytes
      + serializes an HTTP response to wire bytes including status line and headers
      # http_serialization
  std.tcp
    std.tcp.listen
      @ (port: i32) -> result[tcp_listener, string]
      + opens a TCP listener on the given port
      - returns error when the port is in use
      - returns error when the port is privileged and the process lacks capability
      # networking
    std.tcp.accept
      @ (l: tcp_listener) -> result[tcp_conn, string]
      + accepts an inbound connection and returns a connection handle
      # networking

http_server
  http_server.route_match
    @ (routes: list[route], method: string, path: string) -> optional[request_handler]
    + returns the handler for an exact method+path match
    + supports path parameters like /users/{id}
    + returns none when no route matches
    # routing
  http_server.handle_connection
    @ (conn: tcp_conn, routes: list[route]) -> result[void, string]
    + reads the request, routes it, writes the response, and closes the connection
    - returns error when the request cannot be parsed
    # request_handling
    -> std.http.parse_request
    -> std.http.format_response
  http_server.serve
    @ (port: i32, routes: list[route]) -> result[void, string]
    + listens on the port and handles incoming connections in a loop
    + each connection is routed and responded to
    # server_lifecycle
    -> std.tcp.listen
    -> std.tcp.accept
