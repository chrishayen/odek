# Requirement: "a library that generates TOTP and HOTP one-time codes"

Two entry points on top of HMAC-SHA1, base32 decoding, and a time source.

std
  std.encoding
    std.encoding.base32_decode
      @ (text: string) -> result[bytes, string]
      + decodes RFC 4648 base32, accepting optional padding and ignoring case
      - returns error on characters outside the base32 alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha1
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA1, returning 20 bytes
      # cryptography
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

otpgen
  otpgen.hotp
    @ (secret: string, counter: u64, digits: i32) -> result[string, string]
    + returns a zero-padded decimal code of the given length (6 or 8) derived from counter
    - returns error when the secret is not valid base32
    - returns error when digits is outside [6, 8]
    # otp
    -> std.encoding.base32_decode
    -> std.crypto.hmac_sha1
  otpgen.totp
    @ (secret: string, period_seconds: i64, digits: i32) -> result[string, string]
    + derives the counter from current time divided by period_seconds, then delegates to hotp
    - returns error when period_seconds is not positive
    # otp
    -> std.time.now_seconds
