# Requirement: "base64 encode and decode"

Two functions. Both are generic enough to live in std — any project doing binary-to-text transport needs them.

std
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes to base64 text with padding
      + returns "" when given empty bytes
      + the standard "Man" => "TWFu" vector passes
      # encoding
    std.encoding.base64_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes a padded base64 string back to bytes
      + accepts input with or without trailing "=" padding
      - returns error on characters outside the base64 alphabet
      - returns error when length (after padding normalization) is not a multiple of 4
      # encoding

base64
  base64.encode
    fn (data: bytes) -> string
    + encodes bytes to base64 text
    # encoding
    -> std.encoding.base64_encode
  base64.decode
    fn (encoded: string) -> result[bytes, string]
    + decodes base64 text to bytes
    - returns error on invalid input
    # encoding
    -> std.encoding.base64_decode
