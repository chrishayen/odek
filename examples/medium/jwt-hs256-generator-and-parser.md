# Requirement: "a jwt generator and parser (HS256)"

Two project entry points sit atop encoding, hashing, and json primitives.

std
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes to base64url without padding
      + returns "" for empty input
      # encoding
    std.encoding.base64url_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on characters outside the base64url alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      + returns 32 bytes
      # cryptography
    std.crypto.constant_time_eq
      fn (a: bytes, b: bytes) -> bool
      + returns true when a and b are byte-equal, in constant time
      # cryptography
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

sjwt
  sjwt.sign
    fn (claims: map[string, string], secret: string) -> result[string, string]
    + returns a token in "header.payload.signature" form
    - returns error when secret is empty
    # signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  sjwt.parse
    fn (token: string, secret: string) -> result[map[string, string], string]
    + returns the claims when the signature verifies
    - returns error when the token does not have three segments
    - returns error when the signature does not match
    # parsing
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_eq
    -> std.json.parse_object
