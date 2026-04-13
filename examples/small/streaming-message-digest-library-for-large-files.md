# Requirement: "a streaming message digest library for large files"

Computes cryptographic and non-cryptographic digests incrementally so that files larger than memory can be hashed.

std
  std.fs
    std.fs.open_read
      @ (path: string) -> result[file_reader, string]
      + opens the file for sequential reading
      - returns error when the file does not exist
      # filesystem
    std.fs.read_chunk
      @ (r: file_reader, max: u32) -> result[bytes, string]
      + returns up to max bytes from the reader, empty bytes at eof
      # filesystem

streaming_digest
  streaming_digest.new
    @ (algorithm: string) -> result[digest_state, string]
    + creates a digest context for one of md5, sha1, sha256, crc32, blake2s
    - returns error on unknown algorithm name
    # construction
  streaming_digest.update
    @ (state: digest_state, data: bytes) -> digest_state
    + feeds more data into the running digest
    # hashing
  streaming_digest.finalize
    @ (state: digest_state) -> string
    + returns the lowercase hex digest and invalidates the state
    # hashing
  streaming_digest.hash_file
    @ (path: string, algorithm: string) -> result[string, string]
    + returns the hex digest of the entire file, reading it in chunks
    - returns error when the file cannot be opened
    # convenience
    -> std.fs.open_read
    -> std.fs.read_chunk
