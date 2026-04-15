# Requirement: "a peer-to-peer file transfer library between two machines"

Core library exposes a sender and receiver that exchange a file over a network connection using a shared passphrase. Network and crypto primitives live in std.

std
  std.net
    std.net.dial
      fn (host: string, port: u16) -> result[conn_handle, string]
      + returns a connected handle to the remote endpoint
      - returns error when the host is unreachable
      # networking
    std.net.listen
      fn (port: u16) -> result[listener_handle, string]
      + returns a listener bound to the given port
      - returns error when the port is already in use
      # networking
    std.net.accept
      fn (listener: listener_handle) -> result[conn_handle, string]
      + returns the next inbound connection
      # networking
    std.net.send_frame
      fn (conn: conn_handle, data: bytes) -> result[void, string]
      + writes a length-prefixed frame
      # networking
    std.net.recv_frame
      fn (conn: conn_handle) -> result[bytes, string]
      + reads one length-prefixed frame
      - returns error when the peer closes mid-frame
      # networking
  std.crypto
    std.crypto.derive_key
      fn (passphrase: string, salt: bytes) -> bytes
      + deterministically derives a 32-byte key from passphrase and salt
      # cryptography
    std.crypto.seal
      fn (key: bytes, plaintext: bytes) -> bytes
      + authenticated-encrypts plaintext under key
      # cryptography
    std.crypto.open
      fn (key: bytes, ciphertext: bytes) -> result[bytes, string]
      + decrypts and verifies ciphertext
      - returns error when the tag does not validate
      # cryptography

portal
  portal.send_file
    fn (host: string, port: u16, passphrase: string, path: string, contents: bytes) -> result[void, string]
    + transfers file contents to the receiver encrypted under a key derived from passphrase
    - returns error when connection fails
    - returns error when the receiver rejects the handshake
    # file_send
    -> std.net.dial
    -> std.net.send_frame
    -> std.net.recv_frame
    -> std.crypto.derive_key
    -> std.crypto.seal
  portal.receive_file
    fn (port: u16, passphrase: string) -> result[received_file, string]
    + accepts one inbound transfer and returns path and decrypted contents
    - returns error when the passphrase does not match
    # file_receive
    -> std.net.listen
    -> std.net.accept
    -> std.net.recv_frame
    -> std.net.send_frame
    -> std.crypto.derive_key
    -> std.crypto.open
