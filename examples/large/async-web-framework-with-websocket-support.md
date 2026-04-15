# Requirement: "an async web framework with websocket support"

Routes HTTP requests, maintains a middleware chain, and upgrades connections to websocket. std handles HTTP parsing plus the SHA-1 and base64 used by the websocket handshake.

std
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses the request line, headers, and body
      - returns error on malformed input
      # parsing
    std.http.encode_response
      fn (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes a response into wire bytes
      # serialization
  std.crypto
    std.crypto.sha1
      fn (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + standard base64 with padding
      # encoding
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s by sep
      # strings

web
  web.new_app
    fn () -> app_state
    + creates an app with no routes and an empty middleware chain
    # construction
  web.route
    fn (app: app_state, method: string, pattern: string, handler_id: string) -> result[void, string]
    + registers a handler for the method+pattern pair
    - returns error when the exact route is already registered
    # registration
    -> std.strings.split
  web.use
    fn (app: app_state, middleware_id: string) -> void
    + appends a middleware to the chain
    # middleware
  web.handle
    fn (app: app_state, raw: bytes) -> result[bytes, string]
    + parses the request, runs middleware then handler, returns response bytes
    - returns 404 response bytes when no route matches
    - returns error when parsing fails
    # dispatch
    -> std.http.parse_request
    -> std.http.encode_response
  web.ws_upgrade
    fn (headers: map[string,string]) -> result[map[string,string], string]
    + returns the response headers that upgrade the connection
    - returns error when Sec-WebSocket-Key header is missing
    # websocket
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  web.ws_decode_frame
    fn (raw: bytes) -> result[ws_frame, string]
    + decodes a single client-to-server websocket frame
    - returns error on truncated input
    # websocket
  web.ws_encode_frame
    fn (opcode: u8, payload: bytes) -> bytes
    + encodes a server-to-client websocket frame
    # websocket
