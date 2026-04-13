# Requirement: "a clustered in-memory cache with per-item TTL"

Sharded local cache with a pluggable peer-replication hook and individual item expiration.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns the current unix time in milliseconds
      # time
  std.hash
    std.hash.fnv1a_64
      @ (data: bytes) -> u64
      + computes the FNV-1a 64-bit hash of the input
      # hashing

cache
  cache.new
    @ (shard_count: i32, max_entries_per_shard: i32) -> cache_state
    + creates a sharded cache with the given shard count and per-shard capacity
    ? shard_count is rounded up to a power of two
    # construction
  cache.set
    @ (state: cache_state, key: string, value: bytes, ttl_millis: i64) -> cache_state
    + stores the value under the hashed shard with expiration at now + ttl
    + evicts the oldest entry when the shard is full
    # write
    -> std.hash.fnv1a_64
    -> std.time.now_millis
  cache.get
    @ (state: cache_state, key: string) -> optional[bytes]
    + returns the value when present and not expired
    - returns none when the entry is missing
    - returns none when the entry has expired and removes it
    # read
    -> std.hash.fnv1a_64
    -> std.time.now_millis
  cache.delete
    @ (state: cache_state, key: string) -> cache_state
    + removes the entry for the key if present
    # write
    -> std.hash.fnv1a_64
  cache.sweep_expired
    @ (state: cache_state) -> cache_state
    + scans all shards and removes entries whose expiration has passed
    # maintenance
    -> std.time.now_millis
  cache.stats
    @ (state: cache_state) -> cache_stats
    + returns hits, misses, evictions, and total entries
    # observability
  cache.attach_peer_sink
    @ (state: cache_state, sink: fn(cache_event) -> void) -> cache_state
    + registers a callback invoked on every set and delete for cluster replication
    ? event delivery ordering is the caller's responsibility
    # clustering
  cache.apply_peer_event
    @ (state: cache_state, event: cache_event) -> cache_state
    + applies a replication event from a peer without re-publishing it
    - ignores events whose timestamp is older than the local entry
    # clustering
    -> std.hash.fnv1a_64
