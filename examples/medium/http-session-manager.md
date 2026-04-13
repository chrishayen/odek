# Requirement: "a session manager for HTTP servers"

Creates opaque session tokens, stores session data behind a pluggable store, and produces cookie headers.

std
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes as base64url without padding
      # encoding
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

session
  session.new_manager
    @ (ttl_seconds: i64) -> manager_state
    + creates an in-memory session manager with the given session lifetime
    # construction
  session.start
    @ (state: manager_state) -> tuple[manager_state, string]
    + creates a new session and returns its token
    # creation
    -> std.crypto.random_bytes
    -> std.encoding.base64url_encode
    -> std.time.now_seconds
  session.get
    @ (state: manager_state, token: string) -> result[map[string,string], string]
    + returns the data map for a session token
    - returns error when the token is unknown or expired
    # retrieval
    -> std.time.now_seconds
  session.put
    @ (state: manager_state, token: string, key: string, value: string) -> result[manager_state, string]
    + sets a key inside a session
    - returns error when the token is unknown or expired
    # mutation
  session.destroy
    @ (state: manager_state, token: string) -> manager_state
    + removes a session; no-op if the token is unknown
    # destruction
  session.renew
    @ (state: manager_state, token: string) -> result[tuple[manager_state, string], string]
    + rotates to a fresh token preserving data and resets the expiry
    - returns error when the token is unknown or expired
    # rotation
    -> std.crypto.random_bytes
  session.cookie_header
    @ (token: string, name: string, secure: bool) -> string
    + returns a Set-Cookie header value for the token with HttpOnly and SameSite=Lax
    # cookie
  session.gc
    @ (state: manager_state) -> manager_state
    + removes expired sessions
    # housekeeping
    -> std.time.now_seconds
