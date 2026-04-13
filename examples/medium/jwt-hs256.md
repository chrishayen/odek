# Requirement: "a JWT signer and verifier (HS256)"

The project layer is just two entry points; all the real work is std primitives (encoding, HMAC, JSON, time) that any crypto / auth project would reuse.

std
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      + returns "" for empty input
      # encoding
    std.encoding.base64url_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on characters outside the base64url alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      + returns 32 bytes
      # cryptography
    std.crypto.constant_time_eq
      @ (a: bytes, b: bytes) -> bool
      + returns true when two slices are equal in constant time
      + returns false when lengths differ
      # cryptography
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as a JSON object
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

jwt
  jwt.sign
    @ (payload: map[string, string], secret: string) -> result[string, string]
    + returns a JWT in "header.payload.signature" format
    + uses the HS256 algorithm
    - returns error when secret is empty
    # token_signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  jwt.verify
    @ (token: string, secret: string) -> result[map[string, string], string]
    + returns the payload map when signature is valid and token is not expired
    - returns error when the token does not have exactly three segments
    - returns error when the signature does not match
    - returns error when the "exp" claim is in the past
    # token_verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_eq
    -> std.json.parse_object
    -> std.time.now_seconds
