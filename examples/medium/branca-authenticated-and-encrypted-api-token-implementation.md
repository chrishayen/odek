# Requirement: "an implementation of Branca authenticated and encrypted API tokens"

A Branca token is XChaCha20-Poly1305 over (version || timestamp || nonce || payload), base62-encoded. Std owns crypto, encoding, and time primitives.

std
  std.crypto
    std.crypto.xchacha20_poly1305_encrypt
      fn (key: bytes, nonce: bytes, aad: bytes, plaintext: bytes) -> result[bytes, string]
      + returns ciphertext || 16-byte tag
      - returns error when key is not 32 bytes
      - returns error when nonce is not 24 bytes
      # cryptography
    std.crypto.xchacha20_poly1305_decrypt
      fn (key: bytes, nonce: bytes, aad: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns the plaintext when the tag verifies
      - returns error when the tag does not verify
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically-random bytes
      # randomness
  std.encoding
    std.encoding.base62_encode
      fn (data: bytes) -> string
      + returns a base62 representation of the input
      # encoding
    std.encoding.base62_decode
      fn (s: string) -> result[bytes, string]
      + decodes a base62 string
      - returns error on characters outside the base62 alphabet
      # encoding
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

branca
  branca.encode
    fn (key: bytes, payload: bytes) -> result[string, string]
    + returns a base62-encoded Branca token with a fresh nonce and current timestamp
    - returns error when key is not 32 bytes
    # token_creation
    -> std.crypto.random_bytes
    -> std.time.now_seconds
    -> std.crypto.xchacha20_poly1305_encrypt
    -> std.encoding.base62_encode
  branca.decode
    fn (key: bytes, token: string, ttl_seconds: i64) -> result[bytes, string]
    + returns the payload when the token decrypts and is within ttl
    - returns error when the token version byte is not 0xBA
    - returns error when the tag does not verify
    - returns error when (now - timestamp) > ttl_seconds
    # token_verification
    -> std.encoding.base62_decode
    -> std.crypto.xchacha20_poly1305_decrypt
    -> std.time.now_seconds
