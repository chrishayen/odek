# Requirement: "a JSON Web Token library (HS256)"

Signs and verifies JWTs using HMAC-SHA256. The project surface is thin; the real work is in generic std primitives.

std
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes as base64url without padding
      + returns "" on empty input
      # encoding
    std.encoding.base64url_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes base64url with or without trailing padding
      - returns error on characters outside the base64url alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + returns 32-byte HMAC-SHA256 digest
      # cryptography
    std.crypto.constant_time_equal
      fn (a: bytes, b: bytes) -> bool
      + returns true when both inputs are byte-equal
      ? comparison time is independent of position of the first mismatch
      # cryptography
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as a JSON object
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

jwt
  jwt.sign
    fn (claims: map[string, string], secret: string) -> result[string, string]
    + returns a token in "header.payload.signature" form
    - returns error when secret is empty
    # token_signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  jwt.verify
    fn (token: string, secret: string) -> result[map[string, string], string]
    + returns the claims when the signature verifies
    - returns error when the token does not have exactly three segments
    - returns error when the signature does not match
    # token_verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_equal
    -> std.json.parse_object
