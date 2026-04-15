# Requirement: "a UUIDv4 generator returning string or bytes"

Two project entry points over a thin random-bytes primitive.

std
  std.random
    std.random.fill_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness

uuid
  uuid.new_bytes
    fn () -> bytes
    + returns 16 bytes with version 4 and variant RFC 4122 bits set
    # generation
    -> std.random.fill_bytes
  uuid.new_string
    fn () -> string
    + returns the canonical 8-4-4-4-12 hex form of a version-4 uuid
    # formatting
    -> std.random.fill_bytes
