# Requirement: "a compact sortable unique id generator with embedded metadata"

IDs encode a timestamp, a partition byte, and a monotonic counter so that lexicographic order matches creation order.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.encoding
    std.encoding.base32_encode
      fn (data: bytes) -> string
      + encodes bytes using Crockford base32 without padding
      # encoding
    std.encoding.base32_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes Crockford base32, accepting both upper and lower case
      - returns error on characters outside the alphabet
      # encoding

sortable_id
  sortable_id.new_generator
    fn (partition: u8) -> generator_state
    + creates a generator for the given partition with counter zero
    # construction
  sortable_id.next
    fn (state: generator_state) -> tuple[string, generator_state]
    + returns a new id and the advanced generator state
    + ids issued within the same millisecond differ by an incrementing counter
    + ids sort lexicographically in issue order
    # generation
    -> std.time.now_millis
    -> std.encoding.base32_encode
  sortable_id.parse
    fn (id: string) -> result[id_parts, string]
    + returns timestamp_ms, partition, and counter decoded from id
    - returns error when id has the wrong length
    - returns error when id contains invalid base32 characters
    # parsing
    -> std.encoding.base32_decode
  sortable_id.timestamp_of
    fn (id: string) -> result[i64, string]
    + returns the embedded creation timestamp in milliseconds
    - returns error when id is malformed
    # inspection
