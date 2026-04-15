# Requirement: "an HTTP session management and identification library"

Sessions are opaque tokens backed by an in-memory store keyed by random ids. Cookie emission is the caller's job; we return the cookie string.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.rand
    std.rand.bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding

session
  session.new_store
    fn (ttl_seconds: i64) -> session_store
    + creates an empty store where each session expires after ttl_seconds of inactivity
    # construction
  session.create
    fn (store: session_store, user_id: string) -> tuple[string, session_store]
    + returns a new random session id bound to user_id and the updated store
    # creation
    -> std.rand.bytes
    -> std.encoding.hex_encode
    -> std.time.now_seconds
  session.lookup
    fn (store: session_store, session_id: string) -> optional[string]
    + returns the user_id when the session exists and has not expired
    - returns none when the session id is unknown
    - returns none when the session has expired
    # lookup
    -> std.time.now_seconds
  session.touch
    fn (store: session_store, session_id: string) -> session_store
    + resets the last-activity timestamp for the session
    ? no-op when the session does not exist
    # maintenance
    -> std.time.now_seconds
  session.destroy
    fn (store: session_store, session_id: string) -> session_store
    + removes the session from the store
    # teardown
  session.format_cookie
    fn (session_id: string, secure: bool) -> string
    + returns a Set-Cookie header value with HttpOnly and optional Secure flags
    # formatting
