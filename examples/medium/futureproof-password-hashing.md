# Requirement: "a password hashing library with algorithm agility"

Hashes embed an algorithm tag and parameters so verification picks the right primitive and upgrades can be detected.

std
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + returns HMAC-SHA256 of data under key
      # cryptography
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
    std.crypto.constant_time_equal
      @ (a: bytes, b: bytes) -> bool
      + returns true when a and b are equal, in time independent of content
      - returns false when lengths differ
      # cryptography
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + returns standard base64 encoding of bytes
      # encoding
    std.encoding.base64_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes standard base64 into bytes
      - returns error on invalid input
      # encoding

password
  password.hash
    @ (plaintext: string, algorithm: string) -> result[string, string]
    + returns an encoded hash string of the form "$alg$params$salt$digest"
    - returns error when algorithm is unknown
    - returns error when plaintext is empty
    # hashing
    -> std.crypto.random_bytes
    -> std.crypto.hmac_sha256
    -> std.encoding.base64_encode
  password.verify
    @ (plaintext: string, encoded: string) -> result[bool, string]
    + returns true when plaintext recomputes to the stored digest
    - returns false when the digest does not match
    - returns error when encoded has an unrecognized algorithm tag
    # verification
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_equal
    -> std.encoding.base64_decode
  password.needs_rehash
    @ (encoded: string, preferred_algorithm: string) -> bool
    + returns true when the encoded hash uses a different algorithm or weaker parameters
    # upgrade
