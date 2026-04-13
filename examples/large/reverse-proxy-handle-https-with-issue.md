# Requirement: "a reverse proxy that terminates HTTPS and obtains certificates from an ACME authority"

The proxy layer routes requests to upstream backends and holds a certificate cache keyed by hostname; certificates are fetched on demand via an ACME client.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses a request line, headers, and body
      - returns error on malformed start line
      # http_parsing
    std.http.encode_response
      @ (status: u16, headers: map[string, string], body: bytes) -> bytes
      + serializes a response into wire bytes
      # http_encoding
  std.tls
    std.tls.handshake
      @ (client_hello: bytes, cert_chain: bytes, key: bytes) -> result[tls_session, string]
      + completes a server handshake and returns the session state
      - returns error when the client hello is malformed
      # tls
    std.tls.server_name
      @ (client_hello: bytes) -> result[string, string]
      + extracts the SNI hostname from a ClientHello
      - returns error when no SNI extension is present
      # tls
  std.crypto
    std.crypto.generate_key
      @ () -> bytes
      + returns a freshly generated private key
      # cryptography
    std.crypto.sign_csr
      @ (key: bytes, common_name: string) -> bytes
      + returns a signed certificate-signing request for the given CN
      # cryptography
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on malformed input
      # serialization

proxy
  proxy.new
    @ (directory_url: string) -> proxy_state
    + creates a proxy bound to an ACME directory URL with an empty cert cache and empty route table
    # construction
  proxy.add_route
    @ (state: proxy_state, host: string, upstream: string) -> proxy_state
    + registers that requests for the given host should be forwarded to the upstream URL
    - returns unchanged state when host is empty
    # routing
  proxy.lookup_upstream
    @ (state: proxy_state, host: string) -> optional[string]
    + returns the registered upstream for a host
    # routing
  proxy.acme_new_account
    @ (state: proxy_state, contact_email: string) -> result[proxy_state, string]
    + registers a new account with the ACME directory and stores the account key
    - returns error when the directory URL is unreachable
    # acme
    -> std.crypto.generate_key
    -> std.json.parse_object
  proxy.acme_request_cert
    @ (state: proxy_state, host: string) -> result[proxy_state, string]
    + runs an ACME order for the host, completes a challenge, and stores the issued cert
    - returns error when the order or challenge fails
    # acme
    -> std.crypto.sign_csr
    -> std.json.parse_object
  proxy.certificate_for
    @ (state: proxy_state, host: string) -> optional[bytes]
    + returns the cached certificate chain for a host
    # certificate_cache
  proxy.handle_client_hello
    @ (state: proxy_state, hello: bytes) -> result[tls_session, string]
    + selects the certificate by SNI and performs the handshake
    - returns error when no certificate is cached for the requested host
    # tls_termination
    -> std.tls.server_name
    -> std.tls.handshake
  proxy.forward_request
    @ (state: proxy_state, raw: bytes) -> result[bytes, string]
    + parses a plaintext request, looks up the upstream, and returns the serialized response
    - returns a 502 response when no upstream is registered
    # forwarding
    -> std.http.parse_request
    -> std.http.encode_response
