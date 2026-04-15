# Requirement: "a websocket server and client library"

Correct framing, handshake, and message send/receive. Relies on std primitives for TCP, base64, and SHA-1 (used in the handshake).

std
  std.net
    std.net.tcp_listen
      fn (host: string, port: i32) -> result[listener_state, string]
      + binds a TCP listener
      # networking
    std.net.tcp_accept
      fn (listener: listener_state) -> result[conn_state, string]
      + blocks until a client connects
      # networking
    std.net.tcp_connect
      fn (host: string, port: i32) -> result[conn_state, string]
      + opens a TCP connection to host:port
      # networking
    std.net.conn_read
      fn (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes
      # networking
    std.net.conn_write
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes the full buffer
      # networking
  std.crypto
    std.crypto.sha1
      fn (data: bytes) -> bytes
      + computes the SHA-1 hash of data
      + returns 20 bytes
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + standard base64 encoding with padding
      # encoding
    std.encoding.base64_decode
      fn (s: string) -> result[bytes, string]
      + decodes standard base64
      - returns error on invalid characters
      # encoding

websocket
  websocket.server_handshake
    fn (conn: conn_state) -> result[void, string]
    + reads the HTTP upgrade request and writes the accept response
    - returns error when Sec-WebSocket-Key header is missing
    ? computes the accept key as base64(sha1(key + magic))
    # handshake
    -> std.net.conn_read
    -> std.net.conn_write
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  websocket.client_handshake
    fn (conn: conn_state, host: string, path: string) -> result[void, string]
    + sends the upgrade request and validates the server accept response
    - returns error when the server does not respond with 101
    # handshake
    -> std.crypto.random_bytes
    -> std.encoding.base64_encode
  websocket.encode_frame
    fn (opcode: i32, payload: bytes, masked: bool) -> bytes
    + encodes a frame with FIN=1 and the given opcode
    + applies XOR masking when masked=true
    ? extended payload length uses 2 bytes for 126-65535 and 8 bytes for larger
    # framing
    -> std.crypto.random_bytes
  websocket.decode_frame
    fn (raw: bytes) -> result[tuple[i32, bytes, i32], string]
    + returns (opcode, payload, consumed_byte_count)
    - returns error when the buffer does not contain a full frame
    - returns error when a client-to-server frame is not masked
    # framing
  websocket.send_message
    fn (conn: conn_state, opcode: i32, payload: bytes, is_client: bool) -> result[void, string]
    + encodes and writes a full message frame
    # messaging
    -> std.net.conn_write
  websocket.recv_message
    fn (conn: conn_state) -> result[tuple[i32, bytes], string]
    + reads and decodes one message, handling control frames transparently
    - returns error when the connection closes mid-frame
    # messaging
    -> std.net.conn_read
  websocket.close
    fn (conn: conn_state, code: i32, reason: string) -> result[void, string]
    + sends a close frame with the given code and reason
    # lifecycle
