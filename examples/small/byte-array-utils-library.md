# Requirement: "a library of utilities for working with byte arrays"

Small bundle of byte-array helpers: hex and base64 conversion, equality, and concatenation of many buffers.

std: (all units exist)

byte_utils
  byte_utils.to_hex
    fn (data: bytes) -> string
    + returns the lowercase hex representation, two characters per byte
    - empty input returns ""
    # encoding
  byte_utils.from_hex
    fn (hex: string) -> result[bytes, string]
    + parses lowercase or uppercase hex
    - returns error when the length is odd or a non-hex character is present
    # decoding
  byte_utils.to_base64
    fn (data: bytes) -> string
    + returns the standard (padded) base64 encoding
    # encoding
  byte_utils.from_base64
    fn (encoded: string) -> result[bytes, string]
    + decodes standard base64 with or without padding
    - returns error on invalid characters
    # decoding
  byte_utils.concat
    fn (parts: list[bytes]) -> bytes
    + returns a single buffer containing all parts in order
    + returns empty bytes when parts is empty
    # assembly
  byte_utils.equals
    fn (a: bytes, b: bytes) -> bool
    + true when lengths and contents match
    ? compares in constant time so it is safe for secret comparisons
    # comparison
