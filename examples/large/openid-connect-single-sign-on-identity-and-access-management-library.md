# Requirement: "an openid connect single sign-on identity and access management library"

Implements the authorization-code flow: code issuance, token exchange, and access-token verification, plus a minimal role/permission layer.

std
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + returns 32 bytes
      # cryptography
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns 32 bytes
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n bytes from a cryptographic RNG
      # randomness
    std.crypto.bcrypt_hash
      fn (password: string, cost: i32) -> result[string, string]
      + returns a bcrypt-encoded hash
      - returns error when cost is out of range
      # password_hashing
    std.crypto.bcrypt_verify
      fn (password: string, hash: string) -> bool
      + returns true iff password matches the stored hash
      # password_hashing
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + returns base64url without padding
      # encoding
    std.encoding.base64url_decode
      fn (data: string) -> result[bytes, string]
      + decodes base64url input with or without padding
      # encoding
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

sso
  sso.register_user
    fn (username: string, password: string) -> result[user_record, string]
    + returns a user with a bcrypt-hashed password
    - returns error when password is shorter than 8 chars
    # users
    -> std.crypto.bcrypt_hash
  sso.authenticate
    fn (u: user_record, password: string) -> bool
    + returns true when the password verifies
    # users
    -> std.crypto.bcrypt_verify
  sso.register_client
    fn (client_id: string, redirect_uri: string) -> client_record
    + registers an oauth client with the given redirect uri
    # clients
  sso.issue_auth_code
    fn (user: user_record, client: client_record, scope: string) -> auth_code
    + returns a one-time code bound to user, client, and scope
    ? code lifetime is 10 minutes
    # authorization_code
    -> std.crypto.random_bytes
    -> std.encoding.base64url_encode
    -> std.time.now_seconds
  sso.exchange_code
    fn (code: auth_code, client: client_record) -> result[id_token, string]
    + returns a signed id token plus an access token
    - returns error when the code is expired
    - returns error when the client does not match
    # token_exchange
    -> std.time.now_seconds
  sso.sign_id_token
    fn (payload: map[string, string], secret: bytes) -> string
    + returns a JWT-style header.payload.signature
    # token_issuance
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  sso.verify_id_token
    fn (token: string, secret: bytes) -> result[map[string, string], string]
    + returns the claims when signature and expiry check out
    - returns error when the signature does not match
    - returns error when the exp claim has passed
    # token_verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.json.parse_object
    -> std.time.now_seconds
  sso.grant_role
    fn (user: user_record, role: string) -> user_record
    + returns a copy with the role added
    # access_control
  sso.check_permission
    fn (user: user_record, permission: string) -> bool
    + returns true when any of the user's roles grants the permission
    # access_control
