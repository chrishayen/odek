# Requirement: "a protocol server supporting http and websocket on top of an async application interface"

The project layer owns protocol detection and dispatch; std provides the socket and framing primitives.

std
  std.net
    std.net.listen_tcp
      @ (port: i32) -> result[listener, string]
      + returns a listener bound to the given port
      - returns error when the port is already in use
      # networking
    std.net.accept
      @ (l: listener) -> result[connection, string]
      + blocks until a client connects and returns the connection
      # networking
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

protocol_server
  protocol_server.detect_protocol
    @ (first_bytes: bytes) -> string
    + returns "http1", "http2", or "unknown" based on the preface
    # detection
  protocol_server.read_http_request
    @ (conn: connection) -> result[http_request, string]
    + reads a request line plus headers until the blank line
    - returns error on a malformed request line
    # http
  protocol_server.is_websocket_upgrade
    @ (req: http_request) -> bool
    + returns true when Upgrade and Connection headers request websocket
    # detection
  protocol_server.accept_websocket
    @ (conn: connection, req: http_request) -> result[ws_session, string]
    + sends the 101 switching-protocols response with the computed accept key
    - returns error when the request lacks a Sec-WebSocket-Key header
    # websocket
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  protocol_server.read_ws_frame
    @ (session: ws_session) -> result[ws_frame, string]
    + decodes one frame including opcode, fin flag, and unmasked payload
    - returns error on a fragmented frame with an invalid continuation
    # websocket
  protocol_server.dispatch_to_app
    @ (app: app_handler, scope: request_scope, body: bytes) -> result[bytes, string]
    + invokes the application handler with the scope and body and returns its response bytes
    - returns error when the handler itself errors
    # dispatch
