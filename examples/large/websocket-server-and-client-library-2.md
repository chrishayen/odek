# Requirement: "a websocket server and client library"

Handshake, frame parsing, and connection lifecycle. Transport is an opaque TCP connection.

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
    std.net.tcp_dial
      fn (host: string, port: u16) -> result[tcp_conn, string]
      + opens a TCP connection to host:port
      # networking
  std.io
    std.io.read_exact
      fn (conn: tcp_conn, n: i32) -> result[bytes, string]
      + reads exactly n bytes or returns an error
      # io
    std.io.write_all
      fn (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      # io
  std.crypto
    std.crypto.sha1
      fn (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes as standard base64 with padding
      # encoding

websocket
  websocket.compute_accept_key
    fn (client_key: string) -> string
    + returns base64(sha1(client_key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
    # handshake
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  websocket.parse_http_request
    fn (raw: bytes) -> result[http_request, string]
    + parses method, path, and headers
    - returns error on malformed request line
    # parsing
  websocket.server_handshake
    fn (conn: tcp_conn) -> result[ws_conn, string]
    + reads the upgrade request and writes the 101 Switching Protocols response
    - returns error when required headers are missing
    # handshake
    -> std.io.read_exact
    -> std.io.write_all
    -> websocket.parse_http_request
    -> websocket.compute_accept_key
  websocket.client_handshake
    fn (conn: tcp_conn, host: string, path: string) -> result[ws_conn, string]
    + sends the upgrade request with a random Sec-WebSocket-Key and validates the server accept
    # handshake
    -> std.io.write_all
    -> std.io.read_exact
    -> std.crypto.random_bytes
    -> std.encoding.base64_encode
    -> websocket.compute_accept_key
  websocket.encode_frame
    fn (opcode: u8, payload: bytes, mask: bool) -> bytes
    + builds a single FIN frame with the correct length and optional masking
    # framing
    -> std.crypto.random_bytes
  websocket.decode_frame
    fn (conn: tcp_conn) -> result[frame, string]
    + reads a single frame, unmasking payload when the mask bit is set
    - returns error on invalid reserved bits or truncated payload
    # framing
    -> std.io.read_exact
  websocket.send_text
    fn (ws: ws_conn, text: string) -> result[void, string]
    + sends a text frame
    # messaging
    -> websocket.encode_frame
    -> std.io.write_all
  websocket.send_binary
    fn (ws: ws_conn, data: bytes) -> result[void, string]
    + sends a binary frame
    # messaging
    -> websocket.encode_frame
    -> std.io.write_all
  websocket.recv_message
    fn (ws: ws_conn) -> result[ws_message, string]
    + reads frames until a complete message is assembled, transparently handling pings
    - returns error when the peer sends a close frame
    # messaging
    -> websocket.decode_frame
  websocket.close
    fn (ws: ws_conn, code: u16, reason: string) -> result[void, string]
    + sends a close frame and shuts down the underlying connection
    # lifecycle
    -> websocket.encode_frame
    -> std.io.write_all
