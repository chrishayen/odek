# Requirement: "a probabilistic data structures service with persistent storage"

A service that hosts named probabilistic sketches (bloom filter, count-min sketch, hyperloglog) and persists them to disk.

std
  std.hash
    std.hash.murmur3_64
      fn (data: bytes, seed: u32) -> u64
      + returns a 64-bit murmur3 hash of the input
      # hashing
    std.hash.fnv1a_64
      fn (data: bytes) -> u64
      + returns a 64-bit fnv1a hash
      # hashing
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_atomic
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes by renaming a temp file to ensure atomicity
      # filesystem
  std.encoding
    std.encoding.varint_encode
      fn (value: u64) -> bytes
      + encodes an unsigned integer as LEB128 varint
      # encoding
    std.encoding.varint_decode
      fn (data: bytes, offset: i32) -> result[tuple[u64, i32], string]
      + decodes a varint and returns (value, bytes_consumed)
      # encoding

sketches
  sketches.bloom_new
    fn (expected_items: i64, false_positive_rate: f64) -> bloom_state
    + sizes the bit array and hash count to meet the target error rate
    # bloom_construction
  sketches.bloom_add
    fn (state: bloom_state, item: bytes) -> bloom_state
    + sets bits corresponding to k hash positions
    # bloom_insertion
    -> std.hash.murmur3_64
    -> std.hash.fnv1a_64
  sketches.bloom_contains
    fn (state: bloom_state, item: bytes) -> bool
    + returns true when all k bits are set
    - returns false when any bit is unset
    ? false positives possible, false negatives impossible
    # bloom_query
    -> std.hash.murmur3_64
    -> std.hash.fnv1a_64
  sketches.cms_new
    fn (width: i32, depth: i32) -> cms_state
    + allocates a depth-by-width counter matrix
    # cms_construction
  sketches.cms_increment
    fn (state: cms_state, item: bytes, by: i64) -> cms_state
    + increments one counter per row at the hashed column
    # cms_update
    -> std.hash.murmur3_64
  sketches.cms_estimate
    fn (state: cms_state, item: bytes) -> i64
    + returns the minimum counter across rows as the frequency estimate
    # cms_query
    -> std.hash.murmur3_64
  sketches.hll_new
    fn (precision: i32) -> hll_state
    + allocates 2^precision registers
    ? precision between 4 and 16 inclusive
    # hll_construction
  sketches.hll_add
    fn (state: hll_state, item: bytes) -> hll_state
    + updates the register at bucket index with max leading-zero count
    # hll_update
    -> std.hash.murmur3_64
  sketches.hll_cardinality
    fn (state: hll_state) -> i64
    + returns the estimated number of distinct items
    # hll_query
  sketches.serialize
    fn (name: string, kind: string, payload: bytes) -> bytes
    + produces a versioned on-disk binary blob
    # serialization
    -> std.encoding.varint_encode
  sketches.deserialize
    fn (blob: bytes) -> result[tuple[string, string, bytes], string]
    - returns error on version mismatch or truncation
    # deserialization
    -> std.encoding.varint_decode
  sketches.store_save
    fn (dir: string, name: string, blob: bytes) -> result[void, string]
    + writes the sketch blob under dir/name
    # persistence
    -> std.fs.write_atomic
  sketches.store_load
    fn (dir: string, name: string) -> result[bytes, string]
    - returns error when the named sketch is absent
    # persistence
    -> std.fs.read_all
