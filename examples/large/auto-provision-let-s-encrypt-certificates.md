# Requirement: "a library that automatically obtains and renews TLS certificates from an ACME certificate authority"

Runs the ACME HTTP-01 flow, caches certificates, and serves a TLS listener that swaps certificates as they rotate.

std
  std.http
    std.http.client_post
      @ (url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + sends an HTTP POST and returns the response
      - returns error on network or TLS failure
      # http
    std.http.client_get
      @ (url: string, headers: map[string,string]) -> result[http_response, string]
      + sends an HTTP GET
      # http
  std.crypto
    std.crypto.rsa_generate
      @ (bits: i32) -> result[rsa_key_pair, string]
      + generates an RSA key pair
      - returns error when bits is too small
      # cryptography
    std.crypto.ecdsa_generate
      @ (curve: string) -> result[ecdsa_key_pair, string]
      + generates an ECDSA key pair for the named curve
      # cryptography
    std.crypto.sign_jws
      @ (key: signing_key, header: map[string,string], payload: bytes) -> result[string, string]
      + produces a JWS compact serialization
      # cryptography
    std.crypto.csr_build
      @ (key: signing_key, domains: list[string]) -> result[bytes, string]
      + builds a certificate signing request covering the given domains
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes as base64url without padding
      # encoding
  std.json
    std.json.encode_object
      @ (obj: map[string,string]) -> string
      + encodes a string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses JSON into a string map
      - returns error on malformed input
      # serialization
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads file bytes
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes file bytes atomically
      # io
  std.net
    std.net.tls_listen
      @ (host: string, port: u16, cert_provider: cert_provider_fn) -> result[tls_listener, string]
      + serves TLS connections, consulting cert_provider on each handshake
      # networking
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

autocert
  autocert.manager_new
    @ (directory_url: string, email: string, cache_dir: string) -> manager_state
    + creates a manager bound to an ACME directory and a local cache
    # construction
  autocert.account_register
    @ (manager: manager_state) -> result[manager_state, string]
    + registers an account key with the ACME server
    - returns error when the server rejects the registration
    # account
    -> std.crypto.ecdsa_generate
    -> std.crypto.sign_jws
    -> std.http.client_post
    -> std.json.encode_object
    -> std.json.parse_object
  autocert.order_new
    @ (manager: manager_state, domains: list[string]) -> result[order_state, string]
    + creates a new certificate order for the given domains
    # ordering
    -> std.http.client_post
    -> std.json.encode_object
  autocert.challenge_serve
    @ (manager: manager_state, token: string, key_auth: string) -> manager_state
    + records an HTTP-01 challenge response for the token
    # challenge
  autocert.challenge_complete
    @ (manager: manager_state, order: order_state) -> result[order_state, string]
    + notifies the ACME server that challenges are ready and polls until validated
    - returns error when validation fails
    # challenge
    -> std.http.client_post
    -> std.http.client_get
  autocert.finalize
    @ (manager: manager_state, order: order_state, cert_key: signing_key) -> result[bytes, string]
    + submits the CSR and downloads the issued certificate chain
    - returns error when finalization fails
    # issuance
    -> std.crypto.csr_build
    -> std.encoding.base64url_encode
    -> std.http.client_post
    -> std.http.client_get
  autocert.cache_store
    @ (manager: manager_state, domain: string, chain: bytes, key: bytes) -> result[void, string]
    + persists the certificate and private key for a domain
    # caching
    -> std.fs.write_all
  autocert.cache_load
    @ (manager: manager_state, domain: string) -> result[optional[cert_bundle], string]
    + loads a cached certificate if present
    # caching
    -> std.fs.read_all
  autocert.get_certificate
    @ (manager: manager_state, domain: string) -> result[cert_bundle, string]
    + returns a valid certificate for the domain, issuing or renewing if needed
    - returns error when the ACME flow fails
    # dispatch
    -> std.time.now_seconds
  autocert.listen_tls
    @ (manager: manager_state, host: string, port: u16) -> result[tls_listener, string]
    + starts a TLS listener that resolves certificates through the manager
    # serving
    -> std.net.tls_listen
