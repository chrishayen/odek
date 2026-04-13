# Requirement: "a thread-safe in-memory key-value cache tuned for large entry counts"

A sharded map with a fixed byte budget. Entries are stored in bucket rings so reclamation does not walk a global free list. Access is concurrent via per-shard locks.

std
  std.hash
    std.hash.xxh64
      @ (data: bytes) -> u64
      + returns a 64-bit hash suitable for sharding and bucketing
      # hashing
  std.sync
    std.sync.new_mutex
      @ () -> mutex
      + returns a new unlocked mutex
      # concurrency
    std.sync.lock
      @ (m: mutex) -> void
      + blocks until the mutex is acquired
      # concurrency
    std.sync.unlock
      @ (m: mutex) -> void
      + releases the mutex
      # concurrency

cache
  cache.new
    @ (max_bytes: i64, shard_count: i32) -> cache_state
    + returns a cache partitioned into shard_count independent shards
    ? each shard gets max_bytes / shard_count budget
    # construction
    -> std.sync.new_mutex
  cache.set
    @ (state: cache_state, key: bytes, value: bytes) -> void
    + stores the entry, evicting the oldest entries in its shard as needed
    # write
    -> std.hash.xxh64
    -> std.sync.lock
    -> std.sync.unlock
  cache.get
    @ (state: cache_state, key: bytes) -> optional[bytes]
    + returns the stored value when present
    - returns none when the key is absent or was evicted
    # read
    -> std.hash.xxh64
    -> std.sync.lock
    -> std.sync.unlock
  cache.delete
    @ (state: cache_state, key: bytes) -> bool
    + returns true when a value was removed
    - returns false when the key was not present
    # delete
    -> std.hash.xxh64
  cache.stats
    @ (state: cache_state) -> cache_stats
    + returns aggregate entry count, byte usage, hits, misses, and evictions
    # observability
