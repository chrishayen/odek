# Requirement: "a customizable HTTP and HTTPS proxy server library"

Accept client connections, parse requests, run a middleware chain over request and response, forward to the origin, and support HTTPS via CONNECT tunneling with on-the-fly certificate generation.

std
  std.net
    std.net.listen_tcp
      @ (addr: string) -> result[listener, string]
      + binds a TCP listener on the given address
      - returns error when the port cannot be bound
      # network
    std.net.accept
      @ (l: listener) -> result[connection, string]
      + blocks until a client connects and returns the new connection
      # network
    std.net.read
      @ (c: connection, n: i32) -> result[bytes, string]
      + reads up to n bytes from the connection
      # network
    std.net.write
      @ (c: connection, data: bytes) -> result[i32, string]
      + writes bytes to the connection and returns the count written
      # network
    std.net.dial_tcp
      @ (addr: string) -> result[connection, string]
      + opens a TCP connection to the given host:port
      # network
  std.tls
    std.tls.wrap_server
      @ (c: connection, cert: bytes, key: bytes) -> result[connection, string]
      + wraps a plaintext connection in a TLS server-side handshake
      # tls
    std.tls.wrap_client
      @ (c: connection, server_name: string) -> result[connection, string]
      + wraps a plaintext connection in a TLS client handshake
      # tls
  std.http
    std.http.parse_request
      @ (data: bytes) -> result[http_request, string]
      + parses a request line, headers, and body from raw bytes
      - returns error on malformed request line or headers
      # http
    std.http.serialize_request
      @ (req: http_request) -> bytes
      + returns the wire-format bytes for a request
      # http
    std.http.parse_response
      @ (data: bytes) -> result[http_response, string]
      + parses a status line, headers, and body from raw bytes
      # http
    std.http.serialize_response
      @ (resp: http_response) -> bytes
      + returns the wire-format bytes for a response
      # http
  std.crypto
    std.crypto.generate_leaf_cert
      @ (ca_cert: bytes, ca_key: bytes, host: string) -> tuple[bytes, bytes]
      + returns (cert, key) signed by the CA for the given host
      # cryptography

http_proxy
  http_proxy.new
    @ (listen_addr: string, ca_cert: bytes, ca_key: bytes) -> proxy
    + creates a proxy configured to listen on the given address with the given CA for HTTPS interception
    # construction
  http_proxy.use_request
    @ (p: proxy, mw: request_middleware) -> proxy
    + appends a request middleware; middlewares see the request before it is forwarded
    # middleware
  http_proxy.use_response
    @ (p: proxy, mw: response_middleware) -> proxy
    + appends a response middleware; middlewares see the response before it is returned to the client
    # middleware
  http_proxy.start
    @ (p: proxy) -> result[proxy, string]
    + binds the listener and begins accepting connections
    - returns error when the listen address cannot be bound
    # lifecycle
    -> std.net.listen_tcp
  http_proxy.stop
    @ (p: proxy) -> proxy
    + closes the listener and drains in-flight connections
    # lifecycle
  http_proxy.handle_connection
    @ (p: proxy, c: connection) -> void
    + dispatches a client connection to the plain or CONNECT handler based on the first request line
    # dispatch
    -> std.net.read
  http_proxy.handle_plain
    @ (p: proxy, req: http_request) -> result[http_response, string]
    + runs request middlewares, forwards to the origin, and runs response middlewares
    # forwarding
    -> std.net.dial_tcp
    -> std.http.serialize_request
    -> std.http.parse_response
  http_proxy.handle_connect
    @ (p: proxy, c: connection, host: string) -> result[void, string]
    + completes a CONNECT handshake, mints a leaf cert for host, and bridges the decrypted stream through the middleware pipeline
    # tunneling
    -> std.crypto.generate_leaf_cert
    -> std.tls.wrap_server
    -> std.tls.wrap_client
  http_proxy.forward
    @ (req: http_request) -> result[http_response, string]
    + opens a connection to the request's host, writes the request, and parses the response
    # io
    -> std.net.dial_tcp
    -> std.net.write
    -> std.http.serialize_request
    -> std.http.parse_response
