# Requirement: "a generator of concise, unambiguous, URL-safe UUIDs"

Produces 128-bit random ids encoded in a base-57 alphabet that omits visually similar characters.

std
  std.random
    std.random.bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # random

shortuuid
  shortuuid.new
    @ () -> string
    + returns a short id encoded in a 57-character alphabet that excludes 0, O, I, l, and similar look-alikes
    ? underlying entropy is 128 bits of random bytes
    # generation
    -> std.random.bytes
  shortuuid.encode
    @ (raw: bytes) -> string
    + encodes 16 raw bytes as a fixed-length short id string
    # encoding
  shortuuid.decode
    @ (id: string) -> result[bytes, string]
    + decodes a short id back to its 16 raw bytes
    - returns error when id contains characters outside the alphabet
    # decoding
