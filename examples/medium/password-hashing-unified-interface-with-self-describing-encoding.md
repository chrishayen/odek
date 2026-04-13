# Requirement: "a unified interface for multiple password hashing algorithms that produces self-describing encoded hashes"

One project surface: hash and verify. Each algorithm is a pluggable backend keyed by name; the encoded string carries the algorithm identifier so verify can dispatch.

std
  std.crypto
    std.crypto.pbkdf2_sha256
      @ (password: bytes, salt: bytes, iterations: i32, dk_len: i32) -> bytes
      + computes PBKDF2-HMAC-SHA256
      # cryptography
    std.crypto.scrypt
      @ (password: bytes, salt: bytes, n: i32, r: i32, p: i32, dk_len: i32) -> bytes
      + computes scrypt with the given parameters
      # cryptography
  std.random
    std.random.bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      # encoding
    std.encoding.base64_decode
      @ (encoded: string) -> result[bytes, string]
      - returns error on invalid base64
      # encoding

password
  password.register_pbkdf2
    @ (state: registry_state, iterations: i32) -> registry_state
    + registers "pbkdf2-sha256" with the given iteration count
    # registration
  password.register_scrypt
    @ (state: registry_state, n: i32, r: i32, p: i32) -> registry_state
    + registers "scrypt" with the given parameters
    # registration
  password.hash
    @ (state: registry_state, algorithm: string, password: string) -> result[string, string]
    + generates a random salt and returns a self-describing encoded hash ("alg$params$salt$digest")
    - returns error when the algorithm is not registered
    # hashing
    -> std.random.bytes
    -> std.encoding.base64_encode
    -> std.crypto.pbkdf2_sha256
    -> std.crypto.scrypt
  password.verify
    @ (state: registry_state, password: string, encoded: string) -> result[bool, string]
    + parses the algorithm identifier and recomputes the digest in constant time
    - returns error when the encoded string is malformed
    - returns error when the algorithm is not registered
    # verification
    -> std.encoding.base64_decode
    -> std.crypto.pbkdf2_sha256
    -> std.crypto.scrypt
  password.needs_rehash
    @ (state: registry_state, encoded: string) -> result[bool, string]
    + returns true when the encoded hash uses weaker parameters than the current registration
    # maintenance
