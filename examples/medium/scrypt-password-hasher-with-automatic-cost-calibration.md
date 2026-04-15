# Requirement: "a scrypt password hashing library with automatic cost calibration"

Scrypt key derivation with a self-contained hash format and a calibration routine that picks parameters from a target duration.

std
  std.rand
    std.rand.bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.crypto
    std.crypto.pbkdf2_hmac_sha256
      fn (password: bytes, salt: bytes, iterations: i32, key_len: i32) -> bytes
      + derives key_len bytes using PBKDF2 with HMAC-SHA256
      # kdf
    std.crypto.salsa20_8_core
      fn (input: bytes) -> bytes
      + returns the 64-byte salsa20/8 core transform of a 64-byte block
      # primitives
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes to base64 with padding
      # encoding
    std.encoding.base64_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes base64 back to bytes
      - returns error on invalid input
      # encoding
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

scrypt
  scrypt.derive_key
    fn (password: bytes, salt: bytes, n: i32, r: i32, p: i32, key_len: i32) -> result[bytes, string]
    + derives a key using scrypt with parameters n (cost), r (block size), p (parallelism)
    - returns error when n is not a power of 2 or parameters would exceed memory limits
    # kdf
    -> std.crypto.pbkdf2_hmac_sha256
    -> std.crypto.salsa20_8_core
  scrypt.hash_password
    fn (password: string) -> result[string, string]
    + returns a self-contained string "scrypt$n$r$p$salt$hash" with a random salt
    # hashing
    -> std.rand.bytes
    -> std.encoding.base64_encode
  scrypt.verify_password
    fn (password: string, encoded: string) -> result[bool, string]
    + returns true when the password matches the encoded hash
    - returns error on malformed encoded string
    # verification
    -> std.encoding.base64_decode
  scrypt.calibrate
    fn (target_duration_ms: i64) -> tuple[i32, i32, i32]
    + returns (n, r, p) such that hashing a fresh password takes approximately the target duration
    ? starts from a safe baseline and doubles n until the measurement exceeds the target
    # calibration
    -> std.time.now_millis
