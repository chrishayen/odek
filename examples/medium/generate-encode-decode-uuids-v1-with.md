# Requirement: "a uuid v1 generator, encoder, and decoder"

Time-based uuids composed of a 60-bit timestamp, a 14-bit clock sequence, and a 48-bit node identifier. The caller chooses whether the node comes from hardware or a secure random source.

std
  std.time
    std.time.now_100ns_since_gregorian
      @ () -> i64
      + returns the count of 100-nanosecond intervals since 1582-10-15 00:00:00 UTC
      # time
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically secure random bytes
      # cryptography

uuid1
  uuid1.random_node
    @ () -> bytes
    + returns a 6-byte random node id with the multicast bit set
    # node
    -> std.crypto.random_bytes
  uuid1.new_generator
    @ (node: bytes) -> generator_state
    + returns a generator seeded with a 6-byte node identifier and a fresh clock sequence
    - returns an error state when node is not exactly 6 bytes
    # construction
    -> std.crypto.random_bytes
  uuid1.next
    @ (gen: generator_state) -> tuple[generator_state, bytes]
    + returns the next 16-byte uuid and the advanced generator
    + bumps the clock sequence when the timestamp has gone backwards
    # generation
    -> std.time.now_100ns_since_gregorian
  uuid1.encode
    @ (id: bytes) -> string
    + returns the canonical hyphenated lowercase hex representation
    - returns "" when id is not 16 bytes
    # encoding
  uuid1.decode
    @ (text: string) -> result[bytes, string]
    + returns the 16 raw bytes from a canonical or brace-wrapped text form
    - returns error when the text is not a valid uuid form
    # decoding
  uuid1.timestamp
    @ (id: bytes) -> result[i64, string]
    + returns the embedded 100-nanosecond timestamp from a v1 uuid
    - returns error when id is not 16 bytes or the version is not 1
    # inspection
