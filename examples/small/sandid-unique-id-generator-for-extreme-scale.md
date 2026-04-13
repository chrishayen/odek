# Requirement: "a library generating unique identifiers dense enough that every grain of sand could have one"

A 128-bit identifier generator combining time and randomness, plus encoding and comparison helpers.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography

sandid
  sandid.new
    @ () -> bytes
    + returns a 16-byte identifier with a time-based prefix and random suffix
    ? the first 6 bytes encode milliseconds to give sortability
    -> std.time.now_millis
    -> std.crypto.random_bytes
    # generation
  sandid.format
    @ (id: bytes) -> string
    + encodes an identifier as a 32-character lowercase hex string
    - returns "" when input is not exactly 16 bytes
    # rendering
  sandid.parse
    @ (text: string) -> result[bytes, string]
    + decodes a 32-character hex string back into a 16-byte id
    - returns error on wrong length or non-hex character
    # parsing
  sandid.compare
    @ (left: bytes, right: bytes) -> i32
    + returns -1, 0, or 1 comparing two ids byte by byte
    # comparison
