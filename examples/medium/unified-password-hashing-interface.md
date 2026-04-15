# Requirement: "a unified interface across different password hashing algorithms"

One project-level verify function dispatches to algorithm-specific hashers identified by an encoded prefix.

std
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
    std.crypto.bcrypt_hash
      fn (password: string, cost: i32) -> result[string, string]
      + returns a bcrypt-encoded hash string
      - returns error when cost is outside the supported range
      # cryptography
    std.crypto.bcrypt_verify
      fn (password: string, encoded: string) -> result[bool, string]
      + returns true when the password matches the encoded hash
      - returns error when the encoded string is malformed
      # cryptography
    std.crypto.argon2id_hash
      fn (password: string, salt: bytes, memory_kib: i32, iterations: i32) -> string
      + returns an argon2id-encoded hash string
      # cryptography
    std.crypto.argon2id_verify
      fn (password: string, encoded: string) -> result[bool, string]
      + returns true when the password matches
      - returns error when the encoded string is malformed
      # cryptography

passwap
  passwap.hash
    fn (password: string, algorithm: string) -> result[string, string]
    + returns a prefixed encoded hash such as "$2b$..." or "$argon2id$..."
    - returns error when the algorithm name is unknown
    - returns error when the password is empty
    # hashing
    -> std.crypto.random_bytes
    -> std.crypto.bcrypt_hash
    -> std.crypto.argon2id_hash
  passwap.verify
    fn (password: string, encoded: string) -> result[bool, string]
    + returns true when the password matches the encoded hash
    + dispatches to the correct algorithm based on the encoded prefix
    - returns error when the encoded string has no recognized prefix
    # verification
    -> std.crypto.bcrypt_verify
    -> std.crypto.argon2id_verify
  passwap.needs_rehash
    fn (encoded: string, preferred_algorithm: string) -> bool
    + returns true when the encoded hash uses a different algorithm than preferred
    + returns true when the parameters are weaker than current defaults
    # policy
