# Requirement: "a web framework with REST routing and WebSocket support"

Two related surfaces: an HTTP router with method-and-path dispatch, and a WebSocket layer that upgrades HTTP connections and exchanges frames. Transport and framing live in std.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from raw request bytes
      - returns error on malformed request line
      # http
    std.http.format_response
      @ (status: i32, headers: map[string, string], body: bytes) -> bytes
      + serializes a status, headers, and body into wire-format bytes
      # http
  std.crypto
    std.crypto.sha1
      @ (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest
      # cryptography
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + encodes bytes to standard base64 with padding
      # encoding
  std.websocket
    std.websocket.decode_frame
      @ (raw: bytes) -> result[ws_frame, string]
      + returns the next frame from a buffer of websocket bytes
      - returns error on malformed framing
      # websocket
    std.websocket.encode_frame
      @ (opcode: i32, payload: bytes) -> bytes
      + encodes an unmasked server frame with the given opcode and payload
      # websocket

web_framework
  web_framework.new_app
    @ () -> app_state
    + returns an empty application with no routes or sockets
    # construction
  web_framework.add_route
    @ (app: app_state, method: string, path: string, handler_id: string) -> app_state
    + registers a handler for the given method and path
    # routing
  web_framework.match_route
    @ (app: app_state, method: string, path: string) -> optional[string]
    + returns the handler id for the first matching route
    - returns none when no route matches
    # routing
  web_framework.handle_http
    @ (app: app_state, raw: bytes) -> result[bytes, string]
    + parses a request, finds a handler, and returns response bytes
    - returns a 404 response when no route matches
    - returns a 400 response when the request is malformed
    # dispatch
    -> std.http.parse_request
    -> std.http.format_response
  web_framework.compute_accept_key
    @ (client_key: string) -> string
    + returns the base64-encoded SHA-1 of client_key concatenated with the websocket magic string
    # websocket_handshake
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  web_framework.upgrade_connection
    @ (app: app_state, raw: bytes) -> result[bytes, string]
    + returns the 101 switching-protocols response when the request contains a valid upgrade header
    - returns error when the request is not a websocket upgrade
    # websocket_handshake
    -> std.http.parse_request
    -> std.http.format_response
  web_framework.new_socket
    @ (id: string) -> socket_state
    + returns an open socket with an empty read buffer
    # construction
  web_framework.receive
    @ (socket: socket_state, incoming: bytes) -> result[tuple[socket_state, list[string]], string]
    + appends bytes, decodes any complete text frames, and returns their payloads
    - returns error on malformed framing
    # websocket
    -> std.websocket.decode_frame
  web_framework.send_text
    @ (payload: string) -> bytes
    + encodes a text frame ready to send over the wire
    # websocket
    -> std.websocket.encode_frame
