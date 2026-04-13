# Requirement: "a simple hashing library that hashes a buffer, string, or file"

A thin ergonomic facade over standard crypto and filesystem primitives.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the SHA-256 digest of data as 32 bytes
      # cryptography
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      # filesystem
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + encodes bytes to lowercase hexadecimal
      # encoding

hasha
  hasha.hash_bytes
    @ (data: bytes) -> string
    + returns the lowercase hex SHA-256 of the buffer
    + returns a 64-character string
    # hashing
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  hasha.hash_string
    @ (text: string) -> string
    + returns the lowercase hex SHA-256 of the UTF-8 bytes of text
    # hashing
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  hasha.hash_file
    @ (path: string) -> result[string, string]
    + returns the lowercase hex SHA-256 of the file contents
    - returns error when the file cannot be read
    # hashing
    -> std.fs.read_all
    -> std.crypto.sha256
