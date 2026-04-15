# Requirement: "a deduplicating archiver with content-defined chunking, compression, and authenticated encryption"

A repository stores unique content chunks addressed by hash; an archive is a manifest of chunk hashes per file. The project layer composes chunking, compression, encryption, and the manifest.

std
  std.hash
    std.hash.sha256
      fn (data: bytes) -> bytes
      + returns a 32-byte SHA-256 digest
      # cryptography
  std.compress
    std.compress.deflate
      fn (data: bytes) -> bytes
      + compresses bytes using DEFLATE
      # compression
    std.compress.inflate
      fn (data: bytes) -> result[bytes, string]
      + decompresses DEFLATE-compressed bytes
      - returns error on corrupted input
      # compression
  std.crypto
    std.crypto.aead_encrypt
      fn (key: bytes, nonce: bytes, plaintext: bytes, aad: bytes) -> bytes
      + returns ciphertext with an authentication tag appended
      + key must be 32 bytes; nonce must be 12 bytes
      # cryptography
    std.crypto.aead_decrypt
      fn (key: bytes, nonce: bytes, ciphertext: bytes, aad: bytes) -> result[bytes, string]
      + returns plaintext when the tag verifies
      - returns error when the tag does not verify
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hexadecimal
      # encoding

dedup_archive
  dedup_archive.new_repo
    fn (key: bytes) -> result[repo_state, string]
    + creates an empty repository with a 32-byte encryption key
    - returns error when key is not exactly 32 bytes
    # construction
  dedup_archive.chunk
    fn (data: bytes, avg_size: i32) -> list[bytes]
    + splits data into content-defined chunks using a rolling hash
    ? average chunk size is a soft target; boundaries fall where the rolling hash hits a mask
    # chunking
  dedup_archive.put_chunk
    fn (repo: repo_state, chunk: bytes) -> tuple[string, repo_state]
    + hashes, compresses, encrypts, and stores a chunk if not already present
    + returns the hex-encoded hash and the updated repository
    # storage
    -> std.hash.sha256
    -> std.compress.deflate
    -> std.crypto.aead_encrypt
    -> std.crypto.random_bytes
    -> std.encoding.hex_encode
  dedup_archive.get_chunk
    fn (repo: repo_state, hash: string) -> result[bytes, string]
    + returns the original chunk bytes for a stored hash
    - returns error when the hash is unknown
    - returns error when decryption fails
    # retrieval
    -> std.crypto.aead_decrypt
    -> std.compress.inflate
  dedup_archive.add_file
    fn (repo: repo_state, path: string, data: bytes, avg_chunk_size: i32) -> tuple[archive_manifest, repo_state]
    + chunks the file and stores each unique chunk
    + returns an archive manifest entry and the updated repository
    # archiving
  dedup_archive.extract_file
    fn (repo: repo_state, manifest: archive_manifest, path: string) -> result[bytes, string]
    + reassembles the file bytes from manifest chunk hashes
    - returns error when any referenced chunk is missing
    # extraction
  dedup_archive.new_manifest
    fn () -> archive_manifest
    + creates an empty manifest
    # construction
  dedup_archive.manifest_paths
    fn (manifest: archive_manifest) -> list[string]
    + returns the list of archived file paths
    # inspection
  dedup_archive.repo_size
    fn (repo: repo_state) -> i64
    + returns the total size of stored (compressed, encrypted) chunks in bytes
    # inspection
  dedup_archive.unique_chunks
    fn (repo: repo_state) -> i32
    + returns the count of distinct chunks in the repository
    # inspection
