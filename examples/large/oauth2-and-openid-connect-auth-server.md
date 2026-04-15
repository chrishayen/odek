# Requirement: "an authentication and authorization server supporting OAuth2 and OpenID Connect"

Full auth server: user credentials, clients, authorization codes, tokens, and discovery. Crypto and storage isolated in std.

std
  std.crypto
    std.crypto.bcrypt_hash
      fn (password: string, cost: i32) -> result[string, string]
      + returns a bcrypt hash at the given cost factor
      # cryptography
    std.crypto.bcrypt_verify
      fn (password: string, hash: string) -> bool
      + returns true when the password matches the hash
      # cryptography
    std.crypto.rsa_sign_sha256
      fn (private_key_pem: string, data: bytes) -> result[bytes, string]
      - returns error when the key is malformed
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically strong random bytes
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
  std.json
    std.json.encode_object
      fn (obj: map[string, dynamic_value]) -> bytes
      # serialization
    std.json.parse_object
      fn (raw: bytes) -> result[map[string, dynamic_value], string]
      - returns error on invalid JSON
      # serialization
  std.time
    std.time.now_seconds
      fn () -> i64
      # time

auth_server
  auth_server.register_user
    fn (store: auth_store, username: string, password: string) -> result[user_id, string]
    + persists the user with a bcrypt-hashed password
    - returns error when the username is already taken
    # user_management
    -> std.crypto.bcrypt_hash
  auth_server.authenticate
    fn (store: auth_store, username: string, password: string) -> result[user_id, string]
    - returns error "invalid_credentials" on mismatch
    # authentication
    -> std.crypto.bcrypt_verify
  auth_server.register_client
    fn (store: auth_store, name: string, redirect_uris: list[string]) -> result[oauth_client, string]
    + returns client_id and client_secret for a new relying party
    - returns error when redirect_uris is empty
    # client_management
    -> std.crypto.random_bytes
    -> std.encoding.base64url_encode
  auth_server.authorize
    fn (store: auth_store, client_id: string, user: user_id, scopes: list[string], redirect_uri: string) -> result[string, string]
    + returns an authorization code bound to the user, client, and scopes
    - returns error when redirect_uri is not registered for the client
    # authorization
    -> std.crypto.random_bytes
    -> std.encoding.base64url_encode
  auth_server.exchange_code
    fn (store: auth_store, client_id: string, client_secret: string, code: string, redirect_uri: string) -> result[oidc_tokens, string]
    + returns access_token, id_token, and refresh_token
    - returns error when code is unknown, expired, or previously used
    - returns error on client authentication failure
    # token_issuance
    -> std.crypto.random_bytes
    -> std.crypto.rsa_sign_sha256
    -> std.encoding.base64url_encode
    -> std.json.encode_object
    -> std.time.now_seconds
  auth_server.introspect_token
    fn (store: auth_store, token: string) -> result[token_introspection, string]
    + returns active, scope, user, client, and expiry
    - returns inactive=true when the token is unknown or expired
    # introspection
    -> std.time.now_seconds
  auth_server.revoke_token
    fn (store: auth_store, token: string) -> result[void, string]
    + marks the token as inactive
    # revocation
  auth_server.refresh
    fn (store: auth_store, client_id: string, client_secret: string, refresh_token: string) -> result[oidc_tokens, string]
    - returns error when refresh_token is invalid or revoked
    # refresh
    -> std.crypto.random_bytes
    -> std.crypto.rsa_sign_sha256
    -> std.encoding.base64url_encode
    -> std.time.now_seconds
  auth_server.userinfo
    fn (store: auth_store, access_token: string) -> result[map[string, dynamic_value], string]
    - returns error when the token lacks the openid scope
    # userinfo
  auth_server.discovery_document
    fn (issuer: string) -> bytes
    + returns a JSON discovery document listing endpoints and algorithms
    # discovery
    -> std.json.encode_object
  auth_server.jwks
    fn (store: auth_store) -> bytes
    + returns the signing key set in JWK format
    # keys
    -> std.json.encode_object
    -> std.encoding.base64url_encode
