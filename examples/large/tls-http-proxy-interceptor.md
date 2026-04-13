# Requirement: "a tls-capable intercepting http proxy library"

Full feature backend for a man-in-the-middle debugging proxy: listens for client connections, terminates TLS using a generated leaf cert, forwards to the upstream, and exposes flows to registered interceptors.

std
  std.net
    std.net.tcp_listen
      @ (host: string, port: i32) -> result[listener_handle, string]
      + begins listening for TCP connections
      - returns error when the port is in use
      # networking
    std.net.tcp_accept
      @ (listener: listener_handle) -> result[conn_handle, string]
      + blocks until a client connects
      # networking
    std.net.tcp_dial
      @ (host: string, port: i32) -> result[conn_handle, string]
      + opens an outbound TCP connection
      # networking
    std.net.conn_read
      @ (conn: conn_handle, max: i32) -> result[bytes, string]
      + reads up to max bytes from a connection
      # networking
    std.net.conn_write
      @ (conn: conn_handle, data: bytes) -> result[void, string]
      + writes data to a connection
      # networking
  std.tls
    std.tls.server_handshake
      @ (conn: conn_handle, cert_pem: bytes, key_pem: bytes) -> result[conn_handle, string]
      + wraps a conn in TLS, presenting the given cert
      - returns error when the client rejects the cert or handshake fails
      # tls
    std.tls.client_handshake
      @ (conn: conn_handle, server_name: string) -> result[conn_handle, string]
      + performs TLS handshake as client with SNI
      # tls
  std.crypto
    std.crypto.generate_rsa_key
      @ (bits: i32) -> result[bytes, string]
      + returns a PEM-encoded RSA private key
      # cryptography
    std.crypto.sign_leaf_cert
      @ (ca_cert: bytes, ca_key: bytes, host: string, leaf_key: bytes) -> result[bytes, string]
      + issues a leaf certificate for host signed by the given CA
      # cryptography
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses request line, headers, and body
      - returns error on malformed request
      # http
    std.http.parse_response
      @ (raw: bytes) -> result[http_response, string]
      + parses status line, headers, and body
      # http
    std.http.serialize_request
      @ (req: http_request) -> bytes
      + serializes to HTTP/1.1 wire format
      # http
    std.http.serialize_response
      @ (resp: http_response) -> bytes
      + serializes to HTTP/1.1 wire format
      # http

proxy
  proxy.new
    @ (ca_cert: bytes, ca_key: bytes) -> proxy_state
    + creates a proxy configured with the given signing CA
    # construction
  proxy.register_request_interceptor
    @ (state: proxy_state, handler_id: string) -> proxy_state
    + registers a handler that may mutate requests before forwarding
    # extension
  proxy.register_response_interceptor
    @ (state: proxy_state, handler_id: string) -> proxy_state
    + registers a handler that may mutate responses before returning
    # extension
  proxy.start
    @ (state: proxy_state, host: string, port: i32) -> result[proxy_state, string]
    + begins listening for client connections
    - returns error when the address cannot be bound
    # lifecycle
    -> std.net.tcp_listen
  proxy.accept_one
    @ (state: proxy_state) -> result[flow_record, string]
    + accepts one client, runs CONNECT/plain flow, and returns the recorded flow
    # connection_handling
    -> std.net.tcp_accept
    -> std.net.tcp_dial
  proxy.terminate_tls
    @ (state: proxy_state, client_conn: conn_handle, host: string) -> result[conn_handle, string]
    + mints a leaf cert for host and completes the TLS handshake with the client
    - returns error when cert issuance fails
    # tls
    -> std.crypto.generate_rsa_key
    -> std.crypto.sign_leaf_cert
    -> std.tls.server_handshake
  proxy.forward_request
    @ (state: proxy_state, req: http_request, host: string, port: i32) -> result[http_response, string]
    + runs request interceptors, connects upstream, and reads the response
    - returns error when upstream is unreachable
    # forwarding
    -> std.net.tcp_dial
    -> std.tls.client_handshake
    -> std.http.serialize_request
    -> std.http.parse_response
  proxy.apply_response_interceptors
    @ (state: proxy_state, resp: http_response) -> http_response
    + runs each registered response interceptor in order
    # extension
  proxy.record_flow
    @ (state: proxy_state, req: http_request, resp: http_response) -> tuple[proxy_state, flow_record]
    + stores the flow for later inspection and returns it
    # recording
  proxy.list_flows
    @ (state: proxy_state) -> list[flow_record]
    + returns all recorded flows in arrival order
    # query
  proxy.replay_flow
    @ (state: proxy_state, flow_id: string) -> result[http_response, string]
    + re-sends a recorded request and returns the new response
    - returns error when flow_id is unknown
    # replay
    -> std.net.tcp_dial
    -> std.http.serialize_request
    -> std.http.parse_response
