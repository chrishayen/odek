# Requirement: "a JWT login service with pluggable authentication backends"

The service authenticates a credential against one of several registered backends, then issues a signed JWT. Backend plugins are looked up by name.

std
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      + returns 32 bytes
      # cryptography
    std.crypto.bcrypt_verify
      fn (password: string, hash: string) -> bool
      + reports whether the password matches the bcrypt hash
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes as base64url without padding
      # encoding
    std.encoding.base64url_decode
      fn (input: string) -> result[bytes, string]
      + decodes base64url
      - returns error on invalid alphabet
      # encoding
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on invalid input
      # serialization
  std.http
    std.http.post
      fn (url: string, headers: map[string, string], body: bytes) -> result[bytes, string]
      + sends a POST request
      - returns error on network failure
      # http
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns the current unix time in seconds
      # time

login_service
  login_service.new
    fn (signing_secret: string, token_ttl_seconds: i64) -> service_state
    + creates a service bound to a signing secret and token lifetime
    # construction
  login_service.register_backend
    fn (state: service_state, name: string, backend_id: string) -> service_state
    + registers a named authentication backend
    # backend_registration
  login_service.authenticate
    fn (state: service_state, backend_name: string, username: string, credential: string) -> result[map[string, string], string]
    + dispatches to the named backend and returns the resulting user claims
    - returns error when the backend is not registered
    - returns error when the backend rejects the credentials
    # authentication
  login_service.authenticate_htpasswd
    fn (stored_hash: string, password: string) -> bool
    + verifies a password against a bcrypt hash from an htpasswd-style file
    # authentication
    -> std.crypto.bcrypt_verify
  login_service.authenticate_oauth2
    fn (token_url: string, code: string, client_id: string, client_secret: string) -> result[map[string, string], string]
    + exchanges an authorization code for user claims at an OAuth2 endpoint
    - returns error on network failure or error response
    # authentication
    -> std.http.post
    -> std.json.parse_object
  login_service.build_claims
    fn (user: map[string, string], issued_at: i64, ttl: i64) -> map[string, string]
    + returns the claims map with iat and exp fields added
    # claims
  login_service.sign_token
    fn (state: service_state, claims: map[string, string]) -> result[string, string]
    + returns a signed JWT for the given claims
    - returns error when claims cannot be encoded
    # token_signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
    -> std.time.now_seconds
  login_service.verify_token
    fn (state: service_state, token: string) -> result[map[string, string], string]
    + returns the claims when the signature is valid and the token has not expired
    - returns error when the token does not have three segments
    - returns error when the signature does not match
    - returns error when the token has expired
    # token_verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.json.parse_object
    -> std.time.now_seconds
  login_service.login
    fn (state: service_state, backend_name: string, username: string, credential: string) -> result[string, string]
    + authenticates and returns a fresh token
    - returns error on failed authentication
    # login
