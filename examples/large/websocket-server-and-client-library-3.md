# Requirement: "a websocket server and client library"

Server and client sides of the websocket protocol, with handshake, framing, and per-connection state.

std
  std.net
    std.net.tcp_connect
      fn (host: string, port: i32) -> result[conn_handle, string]
      + opens a TCP connection to a remote host
      - returns error when the host is unreachable
      # networking
    std.net.tcp_listen
      fn (host: string, port: i32) -> result[listener_handle, string]
      + binds a TCP listener
      - returns error when the port is in use
      # networking
    std.net.tcp_accept
      fn (lis: listener_handle) -> result[conn_handle, string]
      + accepts the next incoming TCP connection
      # networking
    std.net.read_bytes
      fn (conn: conn_handle, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      - returns error on read failure or EOF
      # networking
    std.net.write_bytes
      fn (conn: conn_handle, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      # networking
  std.crypto
    std.crypto.sha1
      fn (data: bytes) -> bytes
      + returns the 20-byte SHA-1 hash
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes as base64
      # encoding
    std.encoding.base64_decode
      fn (s: string) -> result[bytes, string]
      + decodes a base64 string
      - returns error on invalid characters
      # encoding

websocket
  websocket.compute_accept_key
    fn (client_key: string) -> string
    + computes the Sec-WebSocket-Accept value per the protocol
    # handshake
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  websocket.server_handshake
    fn (conn: conn_handle) -> result[ws_state, string]
    + reads an HTTP upgrade request and writes the upgrade response
    - returns error when the request is not a valid websocket upgrade
    # handshake
    -> std.net.read_bytes
    -> std.net.write_bytes
  websocket.client_connect
    fn (host: string, port: i32, path: string) -> result[ws_state, string]
    + opens a TCP connection and performs the client upgrade handshake
    - returns error when the server rejects the upgrade
    # handshake
    -> std.net.tcp_connect
    -> std.crypto.random_bytes
    -> std.encoding.base64_encode
  websocket.encode_frame
    fn (opcode: i32, payload: bytes, is_client: bool) -> bytes
    + encodes a websocket frame with masking when is_client is true
    + writes the appropriate extended length header for large payloads
    # framing
    -> std.crypto.random_bytes
  websocket.decode_frame
    fn (buf: bytes) -> result[tuple[i32, bytes, i32], string]
    + decodes one frame and returns opcode, payload, and bytes consumed
    - returns error on incomplete or malformed frames
    # framing
  websocket.send_text
    fn (ws: ws_state, message: string) -> result[ws_state, string]
    + sends a text frame over the connection
    # messaging
    -> std.net.write_bytes
  websocket.send_binary
    fn (ws: ws_state, data: bytes) -> result[ws_state, string]
    + sends a binary frame
    # messaging
    -> std.net.write_bytes
  websocket.receive
    fn (ws: ws_state) -> result[tuple[i32, bytes, ws_state], string]
    + reads one complete message, reassembling fragments
    - returns error on protocol violation
    # messaging
    -> std.net.read_bytes
  websocket.send_ping
    fn (ws: ws_state, payload: bytes) -> result[ws_state, string]
    + sends a ping control frame
    # keepalive
  websocket.close
    fn (ws: ws_state, code: i32, reason: string) -> result[void, string]
    + sends a close frame with the given code and reason
    # lifecycle
  websocket.listen
    fn (host: string, port: i32, callback_id: i64) -> result[void, string]
    + accepts connections, performs handshakes, and dispatches each to the callback
    - returns error when the listener cannot be bound
    # serving
    -> std.net.tcp_listen
    -> std.net.tcp_accept
