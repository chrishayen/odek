# Requirement: "an http/2 web server library with configurable routing and static file serving"

A configurable server that loads a site config, serves static files, proxies upstreams, and negotiates HTTP/2. The project layer composes std primitives for TCP, TLS, and filesystem.

std
  std.net
    std.net.tcp_listen
      fn (host: string, port: i32) -> result[listener_state, string]
      + binds a TCP listener
      # networking
    std.net.tcp_accept
      fn (listener: listener_state) -> result[conn_state, string]
      + blocks for a client
      # networking
    std.net.conn_read
      fn (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes
      # networking
    std.net.conn_write
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes the full buffer
      # networking
  std.tls
    std.tls.wrap_server
      fn (conn: conn_state, cert: bytes, key: bytes, alpn: list[string]) -> result[tls_conn, string]
      + performs the TLS handshake as server with ALPN negotiation
      - returns error on handshake failure
      # tls
    std.tls.negotiated_alpn
      fn (tls: tls_conn) -> string
      + returns the negotiated ALPN protocol (e.g. "h2" or "http/1.1")
      # tls
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the full file
      - returns error when the file does not exist
      # filesystem
    std.fs.stat
      fn (path: string) -> result[file_stat, string]
      + returns size, mtime, and is_dir
      # filesystem

webserver
  webserver.parse_config
    fn (raw: string) -> result[server_config, string]
    + parses a server config document into structured host blocks
    - returns error on malformed config
    # config
  webserver.mime_for
    fn (path: string) -> string
    + returns a content-type guess based on file extension
    # static_serving
  webserver.serve_file
    fn (root: string, url_path: string) -> result[http_response, string]
    + reads the file under root and returns a 200 response with proper MIME
    - returns 404 when the file is not found
    - returns 403 when the resolved path escapes the root
    # static_serving
    -> std.fs.read_all
    -> std.fs.stat
  webserver.proxy_pass
    fn (upstream: string, req: http_request) -> result[http_response, string]
    + forwards the request to upstream and returns the response
    - returns 502 on upstream error
    # proxy
  webserver.encode_http1_response
    fn (resp: http_response) -> bytes
    + serializes a response in HTTP/1.1 wire format
    # http1
  webserver.encode_http2_frames
    fn (resp: http_response, stream_id: i32) -> list[bytes]
    + serializes a response as HEADERS and DATA frames for HTTP/2
    # http2
  webserver.handle_request
    fn (config: server_config, req: http_request) -> http_response
    + routes the request via the host block rules to files or proxies
    # routing
  webserver.accept_loop
    fn (config: server_config, host: string, port: i32) -> result[void, string]
    + accepts connections, performs TLS, dispatches via ALPN, and serves
    # server_loop
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> std.net.conn_read
    -> std.net.conn_write
    -> std.tls.wrap_server
    -> std.tls.negotiated_alpn
