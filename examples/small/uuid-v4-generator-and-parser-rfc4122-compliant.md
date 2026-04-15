# Requirement: "a library to generate and parse RFC 4122 compliant version-4 UUIDs"

Two entry points backed by a thin random source.

std
  std.random
    std.random.bytes
      fn (n: i32) -> bytes
      + returns n cryptographically strong random bytes
      # randomness

goid
  goid.new_v4
    fn () -> string
    + returns a 36-character "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx" string
    + the version nibble is 4 and the variant nibble is one of 8, 9, a, or b
    # generation
    -> std.random.bytes
  goid.parse
    fn (text: string) -> result[bytes, string]
    + returns the 16 raw bytes of the uuid
    - returns error when the length is not 36
    - returns error when characters outside hex and dashes are present
    - returns error when the version nibble is not 4
    # parsing
