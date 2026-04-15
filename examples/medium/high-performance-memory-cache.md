# Requirement: "a high performance memory-bound cache"

A cache with a byte-size budget. Admission by sampled LFU frequency, eviction by TinyLFU-style minimum-cost victim selection. Time reads go through a thin std primitive.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.hash
    std.hash.hash64
      fn (data: bytes) -> u64
      + returns a 64-bit non-cryptographic hash
      # hashing

cache
  cache.new
    fn (max_bytes: i64) -> cache_state
    + creates a cache with the given byte budget
    ? uses a count-min sketch for frequency estimation
    # construction
  cache.set
    fn (state: cache_state, key: string, value: bytes, cost: i64) -> cache_state
    + inserts or updates the entry, returning new state
    + rejects the entry when cost exceeds max_bytes
    # insertion
    -> std.hash.hash64
    -> std.time.now_millis
  cache.get
    fn (state: cache_state, key: string) -> tuple[optional[bytes], cache_state]
    + returns the value and incremented frequency when present
    - returns none when the key is absent or expired
    # lookup
    -> std.hash.hash64
    -> std.time.now_millis
  cache.delete
    fn (state: cache_state, key: string) -> cache_state
    + removes the entry when present, returning new state
    # removal
  cache.admit
    fn (state: cache_state, incoming_hash: u64, victim_hash: u64) -> bool
    + returns true when incoming estimated frequency exceeds victim's
    - returns false when victim is hotter than incoming
    # admission_policy
  cache.stats
    fn (state: cache_state) -> tuple[i64, i64, i64]
    + returns (bytes_used, hit_count, miss_count)
    # observability
