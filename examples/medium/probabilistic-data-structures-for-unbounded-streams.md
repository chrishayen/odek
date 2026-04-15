# Requirement: "a set of probabilistic data structures for unbounded streams"

Three streaming structures: Bloom filter for membership, Count-Min sketch for frequency, and HyperLogLog for cardinality.

std
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns a 64-bit FNV-1a hash
      # hashing
    std.hash.murmur3_64
      fn (data: bytes, seed: u64) -> u64
      + returns a 64-bit MurmurHash3 digest with the given seed
      # hashing

probds
  probds.bloom_new
    fn (size_bits: i32, num_hashes: i32) -> bloom_state
    + returns an empty Bloom filter with the given bit array size and hash count
    # bloom
  probds.bloom_add
    fn (bloom: bloom_state, item: bytes) -> bloom_state
    + sets the bits for the given item
    # bloom
    -> std.hash.fnv64
    -> std.hash.murmur3_64
  probds.bloom_contains
    fn (bloom: bloom_state, item: bytes) -> bool
    + returns true when all corresponding bits are set
    - returns false when any bit is unset, guaranteeing no false negatives
    ? false positives are possible by design
    # bloom
    -> std.hash.fnv64
    -> std.hash.murmur3_64
  probds.cms_new
    fn (width: i32, depth: i32) -> cms_state
    + returns an empty Count-Min sketch with depth hash rows of given width
    # count_min
  probds.cms_add
    fn (cms: cms_state, item: bytes, count: i64) -> cms_state
    + increments the counters at each row's hashed column by count
    # count_min
    -> std.hash.murmur3_64
  probds.cms_estimate
    fn (cms: cms_state, item: bytes) -> i64
    + returns the minimum across all rows as the frequency estimate
    ? estimate is always >= true count
    # count_min
    -> std.hash.murmur3_64
  probds.hll_new
    fn (precision: i32) -> hll_state
    + returns an empty HyperLogLog with 2^precision registers
    # hyperloglog
  probds.hll_add
    fn (hll: hll_state, item: bytes) -> hll_state
    + updates the register indexed by the top precision bits with the leading-zero count
    # hyperloglog
    -> std.hash.murmur3_64
  probds.hll_estimate
    fn (hll: hll_state) -> i64
    + returns the estimated number of distinct items added
    # hyperloglog
