# Requirement: "a QUIC transport and HTTP/3 library"

A library implementing QUIC connection state, stream multiplexing, and HTTP/3 framing on top of it. Packets are handed in and out as byte buffers; the host owns UDP I/O.

std
  std.crypto
    std.crypto.aead_seal
      @ (key: bytes, nonce: bytes, plaintext: bytes, aad: bytes) -> bytes
      + encrypts plaintext with AEAD and returns ciphertext with tag
      # cryptography
    std.crypto.aead_open
      @ (key: bytes, nonce: bytes, ciphertext: bytes, aad: bytes) -> result[bytes, string]
      + decrypts and verifies AEAD ciphertext
      - returns error on authentication failure
      # cryptography
    std.crypto.hkdf_expand
      @ (secret: bytes, info: bytes, length: i32) -> bytes
      + HKDF-Expand per RFC 5869
      # cryptography
  std.encoding
    std.encoding.varint_encode
      @ (value: u64) -> bytes
      + encodes an unsigned integer as a QUIC variable-length integer
      # encoding
    std.encoding.varint_decode
      @ (data: bytes) -> result[tuple[u64, i32], string]
      + returns (value, bytes_consumed)
      - returns error when the buffer is too short
      # encoding
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

quic
  quic.new_connection
    @ (role: string, initial_secret: bytes) -> connection_state
    + creates a new connection in the initial state for role "client" or "server"
    # construction
    -> std.crypto.hkdf_expand
  quic.process_datagram
    @ (state: connection_state, datagram: bytes) -> result[connection_state, string]
    + decrypts and processes an incoming UDP datagram containing QUIC packets
    - returns error when the packet header is malformed
    - returns error when AEAD decryption fails
    # packet_processing
    -> std.crypto.aead_open
    -> std.encoding.varint_decode
  quic.poll_datagram
    @ (state: connection_state) -> tuple[optional[bytes], connection_state]
    + returns the next datagram to send, if any
    # packet_processing
    -> std.crypto.aead_seal
    -> std.encoding.varint_encode
  quic.open_stream
    @ (state: connection_state, bidirectional: bool) -> tuple[u64, connection_state]
    + allocates a new stream id and returns it with the updated state
    # stream_management
  quic.write_stream
    @ (state: connection_state, stream_id: u64, data: bytes) -> result[connection_state, string]
    + appends application data to a stream send buffer
    - returns error when the stream is unknown or closed for sending
    # stream_management
  quic.read_stream
    @ (state: connection_state, stream_id: u64) -> tuple[bytes, connection_state]
    + drains and returns buffered received data for a stream
    # stream_management
  quic.tick
    @ (state: connection_state) -> connection_state
    + advances timers for loss detection and congestion control
    # timers
    -> std.time.now_millis

http3
  http3.send_request
    @ (state: connection_state, method: string, path: string, headers: map[string, string]) -> result[tuple[u64, connection_state], string]
    + opens a request stream and writes a HEADERS frame, returning the stream id
    - returns error when the connection is not established
    # request
    -> std.encoding.varint_encode
  http3.read_response
    @ (state: connection_state, stream_id: u64) -> result[tuple[i32, map[string, string], bytes], string]
    + returns (status, headers, body) once the response stream is complete
    - returns error when the stream carries a malformed HTTP/3 frame
    # response
    -> std.encoding.varint_decode
