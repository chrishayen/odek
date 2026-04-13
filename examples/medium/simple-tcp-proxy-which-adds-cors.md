# Requirement: "a TCP proxy that injects CORS headers into proxied HTTP responses"

Parses an incoming http request, forwards it to an upstream, and rewrites the response to add cross-origin headers.

std
  std.net
    std.net.tcp_dial
      @ (host: string, port: i32) -> result[tcp_conn, string]
      + opens a tcp connection to the upstream
      - returns error when the host is unreachable
      # networking
    std.net.tcp_write
      @ (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes bytes to the connection
      # networking
    std.net.tcp_read_all
      @ (conn: tcp_conn) -> result[bytes, string]
      + reads until the connection is closed
      # networking
    std.net.tcp_close
      @ (conn: tcp_conn) -> void
      + closes the connection
      # networking
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses the request line, headers, and body
      - returns error on malformed syntax
      # http
    std.http.encode_request
      @ (req: http_request) -> bytes
      + serializes a request back to wire bytes
      # http
    std.http.parse_response
      @ (raw: bytes) -> result[http_response, string]
      + parses the status line, headers, and body
      - returns error on malformed syntax
      # http
    std.http.encode_response
      @ (resp: http_response) -> bytes
      + serializes a response to wire bytes
      # http

cors_proxy
  cors_proxy.new_config
    @ (upstream_host: string, upstream_port: i32, allowed_origin: string) -> proxy_config
    + creates a config pointing at the upstream and the allowed origin
    # construction
  cors_proxy.inject_headers
    @ (resp: http_response, allowed_origin: string) -> http_response
    + adds access-control-allow-origin and access-control-allow-headers
    + overwrites existing cors headers if already present
    # headers
  cors_proxy.preflight_response
    @ (allowed_origin: string) -> http_response
    + builds a 204 response with the cors headers for options preflights
    # headers
  cors_proxy.handle
    @ (config: proxy_config, raw_request: bytes) -> result[bytes, string]
    + returns the wire bytes of the response with cors headers injected
    + short-circuits to a preflight response when the method is options
    - returns error when the request is malformed or the upstream is unreachable
    # proxying
    -> std.http.parse_request
    -> std.http.encode_request
    -> std.http.parse_response
    -> std.http.encode_response
    -> std.net.tcp_dial
    -> std.net.tcp_write
    -> std.net.tcp_read_all
    -> std.net.tcp_close
