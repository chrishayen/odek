# Requirement: "an in-memory key-value cache for large datasets"

Sharded cache with TTL expiration and size-based eviction, designed for large heap footprints.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + returns 64-bit FNV-1a hash of the input
      # hashing

bigcache
  bigcache.new
    @ (shard_count: i32, max_bytes: i64, default_ttl_ms: i64) -> cache_state
    + creates a cache with the given shard count, byte budget, and default TTL
    ? shard_count must be a power of two to allow mask-based routing
    # construction
  bigcache.set
    @ (state: cache_state, key: string, value: bytes) -> cache_state
    + stores key-value with the default TTL, evicting oldest entries when over budget
    # writes
    -> std.hash.fnv64
    -> std.time.now_millis
  bigcache.set_with_ttl
    @ (state: cache_state, key: string, value: bytes, ttl_ms: i64) -> cache_state
    + stores with an explicit TTL overriding the default
    # writes
    -> std.hash.fnv64
    -> std.time.now_millis
  bigcache.get
    @ (state: cache_state, key: string) -> optional[bytes]
    + returns the cached value when present and unexpired
    - returns empty when the key is absent or expired
    # reads
    -> std.hash.fnv64
    -> std.time.now_millis
  bigcache.delete
    @ (state: cache_state, key: string) -> cache_state
    + removes the entry if present
    # writes
    -> std.hash.fnv64
  bigcache.sweep_expired
    @ (state: cache_state) -> tuple[cache_state, i32]
    + removes expired entries across all shards, returning the count removed
    # maintenance
    -> std.time.now_millis
  bigcache.stats
    @ (state: cache_state) -> cache_stats
    + returns counters for hits, misses, evictions, and current byte usage
    # observability
