# Requirement: "a web session management library with pluggable backends"

Create, load, refresh, and destroy sessions identified by opaque tokens. Storage is delegated to a pluggable backend interface.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically secure random bytes
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes as base64url without padding
      # encoding

sessions
  sessions.new_manager
    fn (ttl_seconds: i64, backend: session_backend) -> manager_state
    + returns a manager that writes and reads through the backend
    # construction
  sessions.create
    fn (state: manager_state, user_id: string) -> result[string, string]
    + returns a fresh opaque session token for the user
    + stores the session with expiry now + ttl_seconds
    - returns error on backend write failure
    # lifecycle
    -> std.crypto.random_bytes
    -> std.encoding.base64url_encode
    -> std.time.now_seconds
  sessions.load
    fn (state: manager_state, token: string) -> result[optional[string], string]
    + returns the user_id when the token is present and unexpired
    - returns none when the token is unknown
    - returns none when the session is expired
    - returns error on backend read failure
    # lifecycle
    -> std.time.now_seconds
  sessions.refresh
    fn (state: manager_state, token: string) -> result[bool, string]
    + extends the session expiry by ttl_seconds and returns true
    - returns false when the token is unknown or expired
    - returns error on backend write failure
    # lifecycle
    -> std.time.now_seconds
  sessions.destroy
    fn (state: manager_state, token: string) -> result[void, string]
    + removes the session from the backend
    + is a no-op when the token is unknown
    - returns error on backend write failure
    # lifecycle
