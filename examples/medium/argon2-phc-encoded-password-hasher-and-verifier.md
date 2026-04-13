# Requirement: "an Argon2 password hashing and verification library producing the standard PHC encoded form"

Hashing returns a PHC-format string containing variant, parameters, salt, and tag; verification parses that string and recomputes the tag.

std
  std.crypto
    std.crypto.argon2id
      @ (password: bytes, salt: bytes, time_cost: i32, memory_kib: i32, parallelism: i32, tag_len: i32) -> bytes
      + returns the raw Argon2id hash tag
      # cryptography
  std.random
    std.random.bytes
      @ (n: i32) -> bytes
      + returns n cryptographically secure random bytes
      # random
  std.encoding
    std.encoding.base64_no_pad_encode
      @ (data: bytes) -> string
      + returns the base64 encoding without trailing padding
      # encoding
    std.encoding.base64_no_pad_decode
      @ (raw: string) -> result[bytes, string]
      + decodes base64 without padding
      - returns error on invalid alphabet
      # encoding

argonpass
  argonpass.hash
    @ (password: string, time_cost: i32, memory_kib: i32, parallelism: i32) -> result[string, string]
    + returns a PHC-encoded hash using a fresh random salt
    - returns error on zero or negative cost parameters
    # hashing
    -> std.random.bytes
    -> std.crypto.argon2id
    -> std.encoding.base64_no_pad_encode
  argonpass.verify
    @ (password: string, encoded: string) -> result[bool, string]
    + returns true when the password hashes to the same tag as the encoded record
    - returns error on malformed encoded string
    # verification
    -> std.encoding.base64_no_pad_decode
    -> std.crypto.argon2id
  argonpass.parse_encoded
    @ (encoded: string) -> result[argon_record, string]
    + splits the PHC string into variant, parameters, salt, and tag
    - returns error when variant is not argon2id
    - returns error on missing parameter fields
    # parsing
    -> std.encoding.base64_no_pad_decode
  argonpass.format_encoded
    @ (record: argon_record) -> string
    + renders a record back to its PHC-encoded form
    # formatting
    -> std.encoding.base64_no_pad_encode
  argonpass.needs_rehash
    @ (encoded: string, target_time: i32, target_memory: i32, target_parallelism: i32) -> result[bool, string]
    + returns true when stored parameters are weaker than the targets
    - returns error on malformed encoded string
    # policy
