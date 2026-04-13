# Requirement: "an HTTP client that can present selectable TLS fingerprints on outbound requests"

Standard HTTP request/response with the additional knob of picking a named TLS fingerprint profile for the handshake.

std
  std.net
    std.net.dial_tcp
      @ (host: string, port: i32) -> result[conn_state, string]
      + opens a TCP connection to host:port
      - returns error on dns, timeout, or refused
      # networking
    std.net.read
      @ (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      # networking
    std.net.write
      @ (conn: conn_state, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      # networking
    std.net.close
      @ (conn: conn_state) -> void
      + closes the connection
      # networking
  std.tls
    std.tls.handshake_with_profile
      @ (conn: conn_state, server_name: string, profile: tls_profile) -> result[tls_state, string]
      + performs a TLS handshake using the cipher suites, extensions, and ordering of profile
      - returns error on handshake failure
      # tls
    std.tls.read
      @ (state: tls_state, max: i32) -> result[bytes, string]
      + reads up to max plaintext bytes
      # tls
    std.tls.write
      @ (state: tls_state, data: bytes) -> result[void, string]
      + writes plaintext bytes
      # tls

http_client
  http_client.new
    @ () -> client_state
    + creates a client with an empty fingerprint registry and default headers
    # construction
  http_client.register_fingerprint
    @ (state: client_state, name: string, profile: tls_profile) -> client_state
    + adds a named TLS fingerprint profile to the client
    # configuration
  http_client.select_fingerprint
    @ (state: client_state, name: string) -> result[client_state, string]
    + marks name as the active profile for subsequent requests
    - returns error when the name is not registered
    # configuration
  http_client.encode_request
    @ (method: string, url: string, headers: map[string,string], body: bytes) -> result[bytes, string]
    + serializes the request into an HTTP/1.1 wire frame
    - returns error on a malformed url
    # protocol
  http_client.decode_response
    @ (frame: bytes) -> result[http_response, string]
    + parses the status line, headers, and body of a response frame
    - returns error on truncated input or invalid headers
    # protocol
  http_client.do
    @ (state: client_state, method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
    + dials, performs the TLS handshake with the active profile, sends the request, and reads the response
    - returns error on any network or TLS failure
    - returns error when no fingerprint is active
    # execution
    -> std.net.dial_tcp
    -> std.tls.handshake_with_profile
    -> std.tls.write
    -> std.tls.read
    -> std.net.close
