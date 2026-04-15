# Requirement: "a peer-to-peer encrypted file transfer library between two endpoints over a realtime datagram channel"

Both ends exchange a key, then one side streams file chunks with per-chunk authenticated encryption. Signaling is caller-provided.

std
  std.fs
    std.fs.read_chunk
      fn (path: string, offset: i64, size: i32) -> result[bytes, string]
      + reads size bytes starting at offset
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_chunk
      fn (path: string, offset: i64, data: bytes) -> result[void, string]
      + writes data at the given offset, creating the file if needed
      # filesystem
    std.fs.file_size
      fn (path: string) -> result[i64, string]
      + returns the file size in bytes
      # filesystem
  std.crypto
    std.crypto.x25519_keypair
      fn () -> tuple[bytes, bytes]
      + returns a fresh (private, public) x25519 keypair
      # cryptography
    std.crypto.x25519_shared
      fn (private: bytes, peer_public: bytes) -> result[bytes, string]
      + derives a 32-byte shared secret
      - returns error when the peer public key is invalid
      # cryptography
    std.crypto.chacha20_poly1305_seal
      fn (key: bytes, nonce: bytes, plaintext: bytes, aad: bytes) -> bytes
      + returns ciphertext with the authentication tag appended
      # cryptography
    std.crypto.chacha20_poly1305_open
      fn (key: bytes, nonce: bytes, ciphertext: bytes, aad: bytes) -> result[bytes, string]
      + returns plaintext when the tag verifies
      - returns error on tag mismatch
      # cryptography
  std.hash
    std.hash.sha256
      fn (data: bytes) -> bytes
      + returns 32-byte SHA-256 digest
      # hashing

file_transfer
  file_transfer.new_session
    fn (chunk_size: i32) -> session_state
    + initializes a transfer session with the given chunk size
    # construction
    -> std.crypto.x25519_keypair
  file_transfer.local_public_key
    fn (state: session_state) -> bytes
    + returns the session public key to be sent over the signaling channel
    # handshake
  file_transfer.complete_handshake
    fn (state: session_state, remote_public: bytes) -> result[session_state, string]
    + derives the shared symmetric key for subsequent chunks
    - returns error on an invalid peer public key
    # handshake
    -> std.crypto.x25519_shared
  file_transfer.begin_send
    fn (state: session_state, path: string) -> result[tuple[session_state, transfer_manifest], string]
    + hashes the file and builds a manifest with size, chunk count, and digest
    - returns error when the file cannot be read
    # send
    -> std.fs.file_size
    -> std.hash.sha256
  file_transfer.next_chunk
    fn (state: session_state, path: string, index: i32) -> result[bytes, string]
    + reads the next plaintext chunk from the source file
    - returns error when index is out of range
    # send
    -> std.fs.read_chunk
  file_transfer.seal_chunk
    fn (state: session_state, index: i32, plaintext: bytes) -> result[bytes, string]
    + encrypts and authenticates the chunk with a nonce derived from index
    - returns error when the session is not yet keyed
    # send
    -> std.crypto.chacha20_poly1305_seal
  file_transfer.open_chunk
    fn (state: session_state, index: i32, ciphertext: bytes) -> result[bytes, string]
    + decrypts a received chunk
    - returns error on authentication failure
    # receive
    -> std.crypto.chacha20_poly1305_open
  file_transfer.write_chunk
    fn (state: session_state, path: string, index: i32, plaintext: bytes) -> result[void, string]
    + writes the decrypted chunk at the correct offset
    - returns error when the destination file cannot be written
    # receive
    -> std.fs.write_chunk
  file_transfer.verify_complete
    fn (state: session_state, path: string, manifest: transfer_manifest) -> result[void, string]
    + rehashes the received file and compares against the manifest digest
    - returns error on digest mismatch
    # verification
    -> std.hash.sha256
