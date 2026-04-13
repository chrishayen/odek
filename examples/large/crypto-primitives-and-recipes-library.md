# Requirement: "a cryptographic primitives and recipes library"

A thin project facade over std crypto primitives. The recipes combine primitives into one-shot "encrypt this message" / "verify this token" helpers.

std
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # hashing
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + returns the 32-byte HMAC-SHA256 tag
      # mac
    std.crypto.aes_gcm_seal
      @ (key: bytes, nonce: bytes, plaintext: bytes, aad: bytes) -> result[bytes, string]
      + returns ciphertext with appended 16-byte tag
      - returns error when key is not 16 or 32 bytes
      # aead
    std.crypto.aes_gcm_open
      @ (key: bytes, nonce: bytes, ciphertext: bytes, aad: bytes) -> result[bytes, string]
      + returns plaintext when the tag verifies
      - returns error on tag mismatch
      # aead
    std.crypto.pbkdf2_sha256
      @ (password: bytes, salt: bytes, iterations: i32, key_len: i32) -> bytes
      + derives a key from a password and salt
      # kdf
    std.crypto.constant_time_equal
      @ (a: bytes, b: bytes) -> bool
      + compares two byte slices without short-circuiting
      # mac
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + standard base64 with padding
      # encoding
    std.encoding.base64_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes standard base64
      - returns error on invalid characters
      # encoding

crypto
  crypto.hash
    @ (data: bytes) -> bytes
    + returns the SHA-256 digest of data
    # hashing
    -> std.crypto.sha256
  crypto.derive_key
    @ (password: string, salt: bytes, iterations: i32) -> bytes
    + derives a 32-byte key suitable for symmetric encryption
    # kdf
    -> std.crypto.pbkdf2_sha256
  crypto.encrypt_message
    @ (key: bytes, plaintext: bytes) -> result[string, string]
    + returns a base64 string containing a fresh nonce and the sealed ciphertext
    - returns error when key length is invalid
    # aead
    -> std.crypto.random_bytes
    -> std.crypto.aes_gcm_seal
    -> std.encoding.base64_encode
  crypto.decrypt_message
    @ (key: bytes, token: string) -> result[bytes, string]
    + returns the original plaintext when the token verifies
    - returns error on tampered or truncated tokens
    # aead
    -> std.encoding.base64_decode
    -> std.crypto.aes_gcm_open
  crypto.sign
    @ (key: bytes, data: bytes) -> bytes
    + returns an HMAC-SHA256 tag for data under key
    # mac
    -> std.crypto.hmac_sha256
  crypto.verify
    @ (key: bytes, data: bytes, tag: bytes) -> bool
    + returns true only when the tag matches in constant time
    # mac
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_equal
