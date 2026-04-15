# Requirement: "a library for handling WebSocket connections"

A full WebSocket implementation: the opening handshake, frame encode/decode with masking, control frame handling (ping/pong/close), and a per-connection state machine.

std
  std.crypto
    std.crypto.sha1
      fn (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes as standard base64 with padding
      # encoding
  std.random
    std.random.bytes
      fn (n: i32) -> bytes
      + returns n random bytes suitable for masking keys
      # randomness
  std.bits
    std.bits.read_u16_be
      fn (data: bytes, offset: i32) -> u16
      + reads a big-endian 16-bit value
      # binary
    std.bits.read_u64_be
      fn (data: bytes, offset: i32) -> u64
      + reads a big-endian 64-bit value
      # binary
    std.bits.write_u16_be
      fn (value: u16) -> bytes
      + encodes a 16-bit value as 2 big-endian bytes
      # binary

websocket
  websocket.accept_key
    fn (client_key: string) -> string
    + returns the Sec-WebSocket-Accept value (base64 of sha1(client_key + magic))
    # handshake
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  websocket.parse_client_handshake
    fn (raw: string) -> result[handshake_request, string]
    + returns method, path, and headers when the HTTP request is a valid upgrade
    - returns error when Upgrade header is not "websocket"
    - returns error when Sec-WebSocket-Key is missing
    # handshake
  websocket.server_handshake_response
    fn (req: handshake_request) -> string
    + builds the 101 Switching Protocols response with the accept key
    # handshake
  websocket.encode_frame
    fn (opcode: i32, payload: bytes, is_client: bool) -> bytes
    + encodes a single frame with fin=1 and masks the payload when is_client is true
    ? extended payload lengths use 2 or 8 bytes per the protocol
    # framing
    -> std.bits.write_u16_be
    -> std.random.bytes
  websocket.decode_frame
    fn (data: bytes) -> result[frame, string]
    + returns the parsed frame and number of bytes consumed
    - returns error "incomplete" when data does not hold a full frame
    - returns error "reserved bits set" when RSV1-3 are non-zero
    # framing
    -> std.bits.read_u16_be
    -> std.bits.read_u64_be
  websocket.new_connection
    fn (is_client: bool) -> conn_state
    + creates a connection in the OPEN state
    # lifecycle
  websocket.feed_bytes
    fn (state: conn_state, data: bytes) -> tuple[conn_state, list[frame]]
    + appends bytes to the read buffer and extracts all complete frames
    # lifecycle
  websocket.send_text
    fn (state: conn_state, payload: string) -> tuple[conn_state, bytes]
    + encodes a text frame and returns the bytes to transmit
    # messaging
  websocket.send_binary
    fn (state: conn_state, payload: bytes) -> tuple[conn_state, bytes]
    + encodes a binary frame
    # messaging
  websocket.handle_control
    fn (state: conn_state, f: frame) -> tuple[conn_state, optional[bytes]]
    + responds to ping with pong, absorbs pong, closes on close
    # control_frames
  websocket.close
    fn (state: conn_state, code: i32, reason: string) -> tuple[conn_state, bytes]
    + transitions to CLOSING and returns a close frame with the given code and reason
    - returns unchanged state and empty bytes when already CLOSED
    # lifecycle
