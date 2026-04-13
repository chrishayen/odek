# Requirement: "a peer-to-peer file and folder transfer library with end-to-end encryption"

A sender hashes a short code and uses it as the PAKE password; a rendezvous server matches sender and receiver; payloads are encrypted with a symmetric cipher derived from the PAKE shared secret. Archiving folders and chunking are handled by the project layer.

std
  std.fs
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns regular file paths under root, recursively
      # filesystem
    std.fs.read_chunk
      @ (path: string, offset: i64, length: i32) -> result[bytes, string]
      + reads a range of bytes from a file
      - returns error on read failure
      # filesystem
    std.fs.write_chunk
      @ (path: string, offset: i64, data: bytes) -> result[void, string]
      + writes bytes at the given offset, extending the file as needed
      # filesystem
  std.archive
    std.archive.pack_tar
      @ (paths: list[string], root: string) -> result[bytes, string]
      + serializes paths relative to root as a tar byte stream
      # archive
    std.archive.unpack_tar
      @ (data: bytes, target_root: string) -> result[list[string], string]
      + extracts a tar byte stream into target_root, returning created paths
      # archive
  std.crypto
    std.crypto.pake_start
      @ (role: i32, password: string) -> pake_state
      + initializes a PAKE state for sender (0) or receiver (1)
      # cryptography
    std.crypto.pake_exchange
      @ (state: pake_state, peer_message: bytes) -> result[pake_step, string]
      + advances the PAKE exchange by consuming a peer message and producing the next
      - returns error on malformed peer message
      # cryptography
    std.crypto.pake_finish
      @ (state: pake_state) -> result[bytes, string]
      + returns the shared symmetric key once the exchange completes
      - returns error when the exchange has not completed
      # cryptography
    std.crypto.aead_seal
      @ (key: bytes, nonce: bytes, plaintext: bytes) -> bytes
      + encrypts and authenticates plaintext under the key and nonce
      # cryptography
    std.crypto.aead_open
      @ (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + decrypts and verifies ciphertext
      - returns error when authentication fails
      # cryptography
  std.net
    std.net.relay_connect
      @ (relay_addr: string, channel_id: string) -> result[relay_conn, string]
      + opens a relay connection under a shared channel identifier
      # network
    std.net.relay_send
      @ (conn: relay_conn, frame: bytes) -> result[void, string]
      + sends a framed message to the peer
      # network
    std.net.relay_recv
      @ (conn: relay_conn) -> result[bytes, string]
      + receives the next framed message from the peer
      # network
    std.net.relay_close
      @ (conn: relay_conn) -> result[void, string]
      + closes the relay connection
      # network
  std.id
    std.id.random_code
      @ (word_count: i32) -> string
      + returns a short human-friendly code of the given word count
      # identity

peerdrop
  peerdrop.new_sender
    @ (relay_addr: string) -> result[sender_state, string]
    + generates a random code and connects to the relay under a channel derived from it
    # sender_setup
    -> std.id.random_code
    -> std.net.relay_connect
  peerdrop.sender_handshake
    @ (state: sender_state) -> result[transport, string]
    + performs the PAKE exchange and returns an encrypted transport
    - returns error when the receiver aborts or the exchange fails
    # handshake
    -> std.crypto.pake_start
    -> std.crypto.pake_exchange
    -> std.crypto.pake_finish
    -> std.net.relay_send
    -> std.net.relay_recv
  peerdrop.new_receiver
    @ (relay_addr: string, code: string) -> result[receiver_state, string]
    + connects to the relay under the channel derived from the code
    # receiver_setup
    -> std.net.relay_connect
  peerdrop.receiver_handshake
    @ (state: receiver_state) -> result[transport, string]
    + performs the PAKE exchange and returns an encrypted transport
    - returns error when the code is wrong and the exchange fails
    # handshake
    -> std.crypto.pake_start
    -> std.crypto.pake_exchange
    -> std.crypto.pake_finish
    -> std.net.relay_send
    -> std.net.relay_recv
  peerdrop.send_path
    @ (transport: transport, path: string) -> result[i64, string]
    + archives the path (file or folder) and streams it as encrypted chunks, returning total bytes sent
    - returns error on read or send failure
    # send
    -> std.fs.walk
    -> std.fs.read_chunk
    -> std.archive.pack_tar
    -> std.crypto.aead_seal
    -> std.net.relay_send
  peerdrop.receive_to
    @ (transport: transport, target_root: string) -> result[list[string], string]
    + receives encrypted chunks, decrypts and reassembles the archive, and extracts it under target_root
    - returns error on authentication or extraction failure
    # receive
    -> std.net.relay_recv
    -> std.crypto.aead_open
    -> std.archive.unpack_tar
    -> std.fs.write_chunk
  peerdrop.close
    @ (transport: transport) -> result[void, string]
    + closes the underlying relay connection
    # teardown
    -> std.net.relay_close
