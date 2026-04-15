# Requirement: "a library tracking users, login state, and permissions using secure cookies and password hashing"

Stores users with hashed passwords, issues a signed session cookie on login, and checks permissions on the current session.

std
  std.crypto
    std.crypto.bcrypt_hash
      fn (password: string, cost: i32) -> string
      + returns a bcrypt hash at the given cost
      # cryptography
    std.crypto.bcrypt_verify
      fn (password: string, hash: string) -> bool
      + reports whether the password matches the hash
      # cryptography
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
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
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

permissions
  permissions.new
    fn (cookie_secret: string) -> store_state
    + creates a store backed by the given secret
    # construction
  permissions.register_user
    fn (state: store_state, username: string, password: string) -> result[store_state, string]
    + adds a user with the given password
    - returns error when the username already exists
    # user_management
    -> std.crypto.bcrypt_hash
  permissions.set_permission
    fn (state: store_state, username: string, permission: string) -> store_state
    + grants a permission to the user
    # permission_management
  permissions.revoke_permission
    fn (state: store_state, username: string, permission: string) -> store_state
    + removes a permission from the user
    # permission_management
  permissions.has_permission
    fn (state: store_state, username: string, permission: string) -> bool
    + reports whether the user holds the permission
    # permission_query
  permissions.verify_password
    fn (state: store_state, username: string, password: string) -> bool
    + reports whether the password matches the stored hash
    # authentication
    -> std.crypto.bcrypt_verify
  permissions.issue_cookie
    fn (state: store_state, username: string, ttl_seconds: i64) -> string
    + returns a signed cookie value binding the username and expiry
    # session
    -> std.time.now_seconds
    -> std.crypto.random_bytes
    -> std.crypto.hmac_sha256
    -> std.encoding.base64url_encode
  permissions.parse_cookie
    fn (state: store_state, cookie: string) -> result[string, string]
    + returns the username when the cookie's signature is valid and it has not expired
    - returns error on signature mismatch
    - returns error when the cookie has expired
    # session
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.time.now_seconds
  permissions.login
    fn (state: store_state, username: string, password: string, ttl_seconds: i64) -> result[string, string]
    + verifies credentials and issues a fresh cookie
    - returns error when credentials are invalid
    # login
  permissions.logout
    fn (state: store_state, cookie: string) -> store_state
    + invalidates the cookie in the store's revocation set
    # logout
