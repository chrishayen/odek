# Requirement: "a library for signing and unsigning data passed through untrusted channels"

Signs a payload with a secret so it can be round-tripped through an untrusted client and verified on return. Optional time-based expiry.

std
  std.encoding
    std.encoding.base64url_encode
      fn (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
    std.encoding.base64url_decode
      fn (encoded: string) -> result[bytes, string]
      - returns error on characters outside the base64url alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha256
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      + returns 32 bytes
      # cryptography
    std.crypto.constant_time_equal
      fn (a: bytes, b: bytes) -> bool
      + returns true when the two byte sequences are equal using constant-time comparison
      # cryptography
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

signer
  signer.sign
    fn (payload: string, secret: string) -> string
    + returns "payload.signature" where signature is HMAC-SHA256(payload) base64url-encoded
    # signing
    -> std.crypto.hmac_sha256
    -> std.encoding.base64url_encode
  signer.unsign
    fn (signed: string, secret: string) -> result[string, string]
    + returns the original payload when the signature matches
    - returns error when the signed value has no '.' separator
    - returns error when the signature does not match
    # verification
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_equal
    -> std.encoding.base64url_decode
  signer.sign_timed
    fn (payload: string, secret: string) -> string
    + returns "payload.timestamp.signature" with the current unix time in the middle segment
    # signing
    -> std.time.now_seconds
    -> std.crypto.hmac_sha256
    -> std.encoding.base64url_encode
  signer.unsign_timed
    fn (signed: string, secret: string, max_age_seconds: i64) -> result[string, string]
    + returns the payload when the signature matches and age is within max_age_seconds
    - returns error when the token is older than max_age_seconds
    - returns error when the signature does not match
    # verification
    -> std.time.now_seconds
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_equal
    -> std.encoding.base64url_decode
