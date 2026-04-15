# Requirement: "an async web framework with routing, middleware, and websocket support"

A framework layer that composes routing, a middleware chain, request/response handling, and websocket upgrades. std carries HTTP parsing, URL decoding, base64, and SHA-1 for the websocket handshake.

std
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses request line, headers, and body
      - returns error on malformed start line or header
      # parsing
    std.http.encode_response
      fn (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes a response into wire bytes
      # serialization
  std.url
    std.url.decode
      fn (s: string) -> result[string, string]
      + percent-decodes a URL component
      - returns error on an incomplete escape sequence
      # parsing
    std.url.parse_query
      fn (raw: string) -> map[string, list[string]]
      + parses a query string into a multi-valued map
      # parsing
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

framework
  framework.new_app
    fn () -> app_state
    + creates an app with no routes and an empty middleware chain
    # construction
  framework.route
    fn (app: app_state, method: string, pattern: string, handler_id: string) -> result[void, string]
    + registers an async handler for the method+pattern
    - returns error when the pattern conflicts with an existing route
    # registration
    -> std.strings.split
  framework.use
    fn (app: app_state, middleware_id: string) -> void
    + appends a middleware to the chain; runs in registration order around the handler
    # middleware
  framework.dispatch
    fn (app: app_state, raw: bytes) -> result[bytes, string]
    + parses the incoming request, runs the middleware chain and handler, returns wire bytes
    - returns error when parsing fails
    - returns 404 wire bytes when no route matches
    # dispatch
    -> std.http.parse_request
    -> std.http.encode_response
    -> std.url.decode
    -> std.url.parse_query
  framework.ws_accept
    fn (headers: map[string,string]) -> result[string, string]
    + computes the websocket accept header value from a client's Sec-WebSocket-Key
    - returns error when the key header is missing
    # websocket
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  framework.ws_decode_frame
    fn (raw: bytes) -> result[ws_frame, string]
    + decodes a single websocket frame, applying the client mask to the payload
    - returns error on truncated input or reserved bit set
    # websocket
  framework.ws_encode_frame
    fn (opcode: u8, payload: bytes) -> bytes
    + encodes a server-to-client frame with no mask
    # websocket
