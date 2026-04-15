# Requirement: "a paranoid library for securely hashing and encrypting passwords"

Derives a key from a password, authenticates, and wraps output in a versioned envelope.

std
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
    std.crypto.scrypt
      fn (password: bytes, salt: bytes, n: i32, r: i32, p: i32, length: i32) -> result[bytes, string]
      + derives a key of given length using scrypt parameters
      - returns error when parameters violate scrypt invariants
      # key_derivation
    std.crypto.hmac_sha512
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA512
      + returns 64 bytes
      # cryptography
    std.crypto.secretbox_seal
      fn (key: bytes, nonce: bytes, plaintext: bytes) -> bytes
      + returns ciphertext with authenticator appended
      # authenticated_encryption
    std.crypto.secretbox_open
      fn (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      - returns error when the authenticator does not verify
      # authenticated_encryption
    std.crypto.constant_time_equal
      fn (a: bytes, b: bytes) -> bool
      + returns true when byte sequences are equal, timing-independent
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + standard base64 with padding
      # encoding
    std.encoding.base64_decode
      fn (encoded: string) -> result[bytes, string]
      - returns error on malformed input
      # encoding

password_vault
  password_vault.hash
    fn (password: string, pepper: bytes) -> result[string, string]
    + returns a versioned envelope "v1$salt$params$mac"
    + uses scrypt to derive an intermediate key then hmac it with the pepper
    - returns error when password is empty
    # password_hashing
    -> std.crypto.random_bytes
    -> std.crypto.scrypt
    -> std.crypto.hmac_sha512
    -> std.encoding.base64_encode
  password_vault.verify
    fn (password: string, pepper: bytes, envelope: string) -> result[bool, string]
    + returns true when envelope matches the password
    - returns error on malformed envelope or unknown version
    # password_verification
    -> std.encoding.base64_decode
    -> std.crypto.scrypt
    -> std.crypto.hmac_sha512
    -> std.crypto.constant_time_equal
  password_vault.encrypt_secret
    fn (password: string, plaintext: bytes) -> result[bytes, string]
    + derives a key from password+random salt, seals plaintext with a random nonce
    + returns salt || nonce || ciphertext
    # authenticated_encryption
    -> std.crypto.random_bytes
    -> std.crypto.scrypt
    -> std.crypto.secretbox_seal
  password_vault.decrypt_secret
    fn (password: string, envelope: bytes) -> result[bytes, string]
    - returns error on wrong password or tampered ciphertext
    # authenticated_decryption
    -> std.crypto.scrypt
    -> std.crypto.secretbox_open
