# Requirement: "a client and server library for encrypted TCP connections"

Symmetric handshake, AEAD-framed messages, and session lifecycle. Cryptographic primitives live in std.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: u16) -> result[tcp_conn, string]
      + returns a connected socket
      - returns error on failure
      # network
    std.net.tcp_listen
      @ (host: string, port: u16) -> result[tcp_listener, string]
      + returns a listener bound to the address
      - returns error when the port is in use
      # network
    std.net.tcp_accept
      @ (listener: tcp_listener) -> result[tcp_conn, string]
      + returns the next accepted connection
      - returns error when the listener is closed
      # network
    std.net.tcp_write
      @ (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes
      # network
    std.net.tcp_read_exact
      @ (conn: tcp_conn, n: i32) -> result[bytes, string]
      + reads exactly n bytes
      - returns error on premature eof
      # network
  std.crypto
    std.crypto.x25519_keypair
      @ () -> tuple[bytes, bytes]
      + returns (private_key, public_key)
      # cryptography
    std.crypto.x25519_shared
      @ (private_key: bytes, peer_public: bytes) -> bytes
      + returns the 32-byte shared secret
      # cryptography
    std.crypto.hkdf_sha256
      @ (secret: bytes, info: bytes, length: i32) -> bytes
      + expands to length bytes
      # cryptography
    std.crypto.aead_seal
      @ (key: bytes, nonce: bytes, plaintext: bytes, aad: bytes) -> bytes
      + returns ciphertext with appended authentication tag
      # cryptography
    std.crypto.aead_open
      @ (key: bytes, nonce: bytes, ciphertext: bytes, aad: bytes) -> result[bytes, string]
      + returns plaintext when the tag verifies
      - returns error on authentication failure
      # cryptography
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography

secure_link
  secure_link.client_handshake
    @ (host: string, port: u16) -> result[secure_session, string]
    + performs key exchange and returns an authenticated session
    - returns error when the peer rejects the handshake
    # handshake
    -> std.net.tcp_connect
    -> std.crypto.x25519_keypair
    -> std.crypto.x25519_shared
    -> std.crypto.hkdf_sha256
  secure_link.server_handshake
    @ (conn: tcp_conn) -> result[secure_session, string]
    + completes the server side of key exchange
    - returns error on protocol violation
    # handshake
    -> std.crypto.x25519_keypair
    -> std.crypto.x25519_shared
    -> std.crypto.hkdf_sha256
  secure_link.listen
    @ (host: string, port: u16) -> result[tcp_listener, string]
    + binds a listener that will accept encrypted connections
    # server
    -> std.net.tcp_listen
  secure_link.accept
    @ (listener: tcp_listener) -> result[secure_session, string]
    + accepts a connection and completes the handshake
    # server
    -> std.net.tcp_accept
  secure_link.send
    @ (session: secure_session, message: bytes) -> result[secure_session, string]
    + encrypts and frames the message, returning the session with advanced nonce
    - returns error when the socket is closed
    # messaging
    -> std.crypto.aead_seal
    -> std.crypto.random_bytes
    -> std.net.tcp_write
  secure_link.recv
    @ (session: secure_session) -> result[tuple[bytes, secure_session], string]
    + decrypts the next framed message
    - returns error on tag mismatch
    # messaging
    -> std.net.tcp_read_exact
    -> std.crypto.aead_open
  secure_link.close
    @ (session: secure_session) -> void
    + closes the underlying connection
    # session
  secure_link.encode_frame
    @ (ciphertext: bytes) -> bytes
    + prepends a 4-byte big-endian length
    # framing
  secure_link.decode_frame
    @ (conn: tcp_conn) -> result[bytes, string]
    + reads a length prefix then that many bytes
    - returns error on eof mid-frame
    # framing
    -> std.net.tcp_read_exact
