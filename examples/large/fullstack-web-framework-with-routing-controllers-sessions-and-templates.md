# Requirement: "a full-stack web framework with routing, controllers, sessions, and templates"

Convention-driven web framework. Transport is an opaque listener; templates use a simple interpolation engine.

std
  std.net
    std.net.tcp_listen
      fn (addr: string, port: u16) -> result[tcp_listener, string]
      + binds a TCP listener
      # networking
    std.net.tcp_accept
      fn (l: tcp_listener) -> result[tcp_conn, string]
      + returns the next accepted connection
      # networking
  std.io
    std.io.read_all
      fn (conn: tcp_conn) -> result[bytes, string]
      + reads until the peer closes or an error occurs
      # io
    std.io.write_all
      fn (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes
      # io
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + returns a 32-byte HMAC
      # cryptography

webapp
  webapp.new_app
    fn (secret: string) -> app_state
    + creates an application with empty routes, no controllers, and an HMAC session key
    # construction
  webapp.route
    fn (app: app_state, method: string, pattern: string, action: fn(request, session) -> response) -> app_state
    + registers an action under method and path pattern
    + pattern segments starting with ":" bind to request params
    # routing
  webapp.match
    fn (app: app_state, method: string, path: string) -> optional[route_match]
    + returns the matched action and extracted params
    # routing
  webapp.parse_request
    fn (raw: bytes) -> result[request, string]
    + parses request line, headers, cookies, and form body
    - returns error on malformed start line
    # parsing
  webapp.encode_response
    fn (resp: response) -> bytes
    + serializes status, headers, and body
    # serialization
  webapp.new_session
    fn () -> session
    + creates an empty session
    # sessions
  webapp.session_put
    fn (s: session, key: string, value: string) -> session
    + sets a key
    # sessions
  webapp.session_get
    fn (s: session, key: string) -> optional[string]
    + returns the value if present
    # sessions
  webapp.encode_session_cookie
    fn (s: session, secret: string) -> string
    + serializes the session and appends an HMAC tag
    # sessions
    -> std.crypto.hmac_sha256
  webapp.decode_session_cookie
    fn (cookie: string, secret: string) -> result[session, string]
    + verifies the HMAC tag and returns the session contents
    - returns error on tag mismatch
    # sessions
    -> std.crypto.hmac_sha256
  webapp.render_template
    fn (source: string, context: map[string, string]) -> string
    + replaces "{{ name }}" placeholders from context, HTML-escaping each value
    # templating
  webapp.handle_connection
    fn (app: app_state, conn: tcp_conn) -> void
    + reads a request, decodes the session cookie, dispatches to a matched action, and writes the response with an updated session cookie
    # request_handling
    -> std.io.read_all
    -> std.io.write_all
    -> webapp.parse_request
    -> webapp.match
    -> webapp.decode_session_cookie
    -> webapp.encode_session_cookie
    -> webapp.encode_response
  webapp.serve
    fn (app: app_state, addr: string, port: u16) -> result[void, string]
    + accepts connections in a loop and dispatches each to handle_connection
    # server
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> webapp.handle_connection
