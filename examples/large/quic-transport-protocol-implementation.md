# Requirement: "an implementation of the IETF QUIC transport protocol"

QUIC is a UDP-based multiplexed transport with TLS 1.3 handshake, congestion control, and stream flow control. Std owns TLS, crypto, UDP, and time primitives.

std
  std.net
    std.net.udp_open
      fn (addr: string) -> result[udp_socket, string]
      + returns a bound UDP socket on the given local address
      - returns error when the address cannot be bound
      # network
    std.net.udp_send
      fn (sock: udp_socket, peer: string, data: bytes) -> result[void, string]
      + sends a datagram to peer
      # network
    std.net.udp_recv
      fn (sock: udp_socket) -> result[tuple[string, bytes], string]
      + returns (peer_addr, datagram) from the next received packet
      # network
  std.crypto
    std.crypto.aes_gcm_encrypt
      fn (key: bytes, nonce: bytes, aad: bytes, plaintext: bytes) -> result[bytes, string]
      + returns ciphertext || tag
      # cryptography
    std.crypto.aes_gcm_decrypt
      fn (key: bytes, nonce: bytes, aad: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns plaintext when the tag verifies
      - returns error when the tag does not verify
      # cryptography
    std.crypto.hkdf_expand_label
      fn (secret: bytes, label: string, length: i32) -> bytes
      + returns derived key material per the TLS 1.3 hkdf-expand-label construction
      # key_derivation
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns cryptographically-random bytes
      # randomness
  std.tls
    std.tls.handshake_client
      fn (server_name: string, initial_data: bytes) -> result[tls_handshake_state, string]
      + returns a handshake state after processing server's initial flight
      - returns error on unsupported cipher suites
      # tls
    std.tls.derive_traffic_keys
      fn (state: tls_handshake_state) -> tuple[bytes, bytes]
      + returns (client_key, server_key) for 1-RTT traffic
      # tls
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

quic
  quic.connect
    fn (peer: string, server_name: string) -> result[quic_connection, string]
    + returns a connection after completing the QUIC handshake
    - returns error when peer is unreachable
    - returns error when the handshake fails
    # connection
    -> std.net.udp_open
    -> std.crypto.random_bytes
    -> std.tls.handshake_client
    -> std.tls.derive_traffic_keys
  quic.accept
    fn (local: string) -> result[quic_listener, string]
    + returns a listener bound to the local address
    # connection
    -> std.net.udp_open
  quic.listener_accept
    fn (l: quic_listener) -> result[quic_connection, string]
    + returns the next incoming connection after the handshake completes
    # connection
  quic.open_stream
    fn (c: quic_connection, bidirectional: bool) -> result[tuple[quic_stream_id, quic_connection], string]
    + returns a new stream id and updated connection state
    - returns error when the peer's stream limit is reached
    # streams
  quic.stream_write
    fn (c: quic_connection, sid: quic_stream_id, data: bytes) -> result[quic_connection, string]
    + buffers data into the stream's send queue, respecting flow control
    - returns error when the stream is already closed for sending
    # streams
  quic.stream_read
    fn (c: quic_connection, sid: quic_stream_id) -> result[tuple[bytes, quic_connection], string]
    + returns any available stream bytes and updated state
    # streams
  quic.encode_frame
    fn (frame: quic_frame) -> bytes
    + returns the on-the-wire encoding of a QUIC frame
    # framing
  quic.decode_frame
    fn (data: bytes) -> result[tuple[quic_frame, i32], string]
    + returns (frame, bytes_consumed)
    - returns error on truncated or unknown frame type
    # framing
  quic.protect_packet
    fn (header: bytes, payload: bytes, key: bytes, packet_number: i64) -> result[bytes, string]
    + returns an encrypted QUIC packet using the appropriate AEAD
    # packet_protection
    -> std.crypto.hkdf_expand_label
    -> std.crypto.aes_gcm_encrypt
  quic.unprotect_packet
    fn (packet: bytes, key: bytes) -> result[tuple[bytes, bytes], string]
    + returns (header, payload) when AEAD verification succeeds
    - returns error when decryption fails
    # packet_protection
    -> std.crypto.aes_gcm_decrypt
  quic.update_loss_detection
    fn (c: quic_connection) -> quic_connection
    + advances the loss detection and congestion control state using current time
    + schedules retransmissions for any declared-lost packets
    # loss_recovery
    -> std.time.now_millis
  quic.close
    fn (c: quic_connection, error_code: i64, reason: string) -> result[void, string]
    + sends a CONNECTION_CLOSE frame and tears the connection down
    # connection
    -> std.net.udp_send
