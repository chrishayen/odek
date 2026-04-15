# Requirement: "a modular hash and cryptography library"

A pluggable collection of common hash and symmetric primitives. Constructions are kept as thin project wrappers that select an algorithm and delegate to std primitives.

std
  std.hash
    std.hash.md5
      fn (data: bytes) -> bytes
      + returns the 16-byte MD5 digest
      # hashing
    std.hash.sha1
      fn (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest
      # hashing
    std.hash.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # hashing
    std.hash.sha512
      fn (data: bytes) -> bytes
      + returns the 64-byte SHA-512 digest
      # hashing
    std.hash.blake2b
      fn (data: bytes, out_len: i32) -> bytes
      + returns a BLAKE2b digest of the requested length
      - returns empty on out_len outside [1, 64]
      # hashing
  std.crypto
    std.crypto.aes_encrypt_block
      fn (key: bytes, block: bytes) -> result[bytes, string]
      + encrypts a 16-byte block under the AES key
      - returns error when block is not 16 bytes or key is not 16/24/32 bytes
      # cryptography
    std.crypto.aes_decrypt_block
      fn (key: bytes, block: bytes) -> result[bytes, string]
      + decrypts a 16-byte AES block
      - returns error on malformed inputs
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness

crypto
  crypto.hash
    fn (algorithm: string, data: bytes) -> result[bytes, string]
    + dispatches to the named digest and returns the raw hash
    - returns error when algorithm is unknown
    # hashing
    -> std.hash.md5
    -> std.hash.sha1
    -> std.hash.sha256
    -> std.hash.sha512
    -> std.hash.blake2b
  crypto.hmac
    fn (algorithm: string, key: bytes, data: bytes) -> result[bytes, string]
    + returns HMAC using the named inner hash
    - returns error when algorithm is unknown
    # mac
    -> std.hash.sha256
    -> std.hash.sha512
  crypto.aes_cbc_encrypt
    fn (key: bytes, iv: bytes, plaintext: bytes) -> result[bytes, string]
    + PKCS#7-pads then encrypts in CBC mode
    - returns error when iv is not 16 bytes
    # cryptography
    -> std.crypto.aes_encrypt_block
  crypto.aes_cbc_decrypt
    fn (key: bytes, iv: bytes, ciphertext: bytes) -> result[bytes, string]
    + decrypts CBC and strips PKCS#7 padding
    - returns error on bad padding
    # cryptography
    -> std.crypto.aes_decrypt_block
  crypto.pbkdf2
    fn (algorithm: string, password: bytes, salt: bytes, iters: i32, key_len: i32) -> result[bytes, string]
    + derives a key via PBKDF2 using the named HMAC
    - returns error when iters < 1 or key_len < 1
    # key_derivation
    -> std.hash.sha256
  crypto.constant_time_equal
    fn (a: bytes, b: bytes) -> bool
    + returns true when the two byte slices match
    ? runs in time proportional to max(len(a), len(b))
    # cryptography
