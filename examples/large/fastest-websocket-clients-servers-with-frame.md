# Requirement: "a websocket client and server with frame-level interface"

Frame-level API so callers can drive the wire protocol directly. Handshake, framing, and masking live in the project package; TCP and crypto are std primitives.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: i32) -> result[tcp_conn, string]
      + opens a TCP connection
      - returns error on unresolved host
      # networking
    std.net.tcp_listen
      @ (host: string, port: i32) -> result[tcp_listener, string]
      + binds and listens on the given address
      # networking
    std.net.tcp_accept
      @ (listener: tcp_listener) -> result[tcp_conn, string]
      + accepts the next incoming connection
      # networking
    std.net.tcp_read
      @ (conn: tcp_conn, n: i32) -> result[bytes, string]
      + reads up to n bytes from the connection
      - returns error when the connection is closed
      # networking
    std.net.tcp_write
      @ (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes
      # networking
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + encodes bytes to standard base64
      # encoding
  std.crypto
    std.crypto.sha1
      @ (data: bytes) -> bytes
      + computes SHA-1 digest
      # cryptography
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns cryptographically random bytes
      # cryptography
  std.http
    std.http.read_request
      @ (conn: tcp_conn) -> result[http_request, string]
      + reads an HTTP/1.1 request line and headers
      - returns error on malformed header
      # parsing
    std.http.write_response
      @ (conn: tcp_conn, status: i32, headers: map[string, string]) -> result[void, string]
      + writes a status line and header block
      # serialization

websocket
  websocket.client_handshake
    @ (host: string, port: i32, path: string) -> result[ws_conn, string]
    + connects, sends the upgrade request, and validates Sec-WebSocket-Accept
    - returns error when the server response is not 101
    # handshake
    -> std.net.tcp_connect
    -> std.net.tcp_write
    -> std.net.tcp_read
    -> std.crypto.random_bytes
    -> std.encoding.base64_encode
  websocket.server_handshake
    @ (conn: tcp_conn) -> result[ws_conn, string]
    + reads the upgrade request and replies with the accept token
    - returns error when the request is missing Sec-WebSocket-Key
    # handshake
    -> std.http.read_request
    -> std.http.write_response
    -> std.crypto.sha1
    -> std.encoding.base64_encode
  websocket.encode_frame
    @ (opcode: u8, payload: bytes, mask: bool) -> bytes
    + produces a framed packet with the chosen opcode
    + masks payload with a random key when mask is true
    ? uses 7/16/64-bit length encoding based on payload size
    # framing
    -> std.crypto.random_bytes
  websocket.decode_frame
    @ (data: bytes) -> result[ws_frame, string]
    + returns opcode, payload, and whether the FIN bit is set
    - returns error when the buffer is shorter than the declared length
    # framing
  websocket.read_frame
    @ (conn: ws_conn) -> result[ws_frame, string]
    + reads a full frame from the connection
    - returns error on short read
    # io
    -> std.net.tcp_read
  websocket.write_frame
    @ (conn: ws_conn, frame: ws_frame) -> result[void, string]
    + writes an encoded frame to the connection
    # io
    -> std.net.tcp_write
  websocket.close
    @ (conn: ws_conn, code: u16, reason: string) -> result[void, string]
    + sends a close frame and releases the connection
    # lifecycle
