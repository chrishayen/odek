# Requirement: "a thread-safe cache with high performance and automatic expiry pruning"

A concurrent in-memory key-value cache with per-entry TTLs and a background pruner. Time and mutexes are thin std primitives.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.sync
    std.sync.new_mutex
      @ () -> mutex_handle
      + creates an unlocked mutex
      # concurrency
    std.sync.with_lock
      @ (m: mutex_handle, body: fn() -> void) -> void
      + runs body while holding the mutex
      # concurrency

gocache
  gocache.new
    @ (default_ttl_ms: i64) -> cache_state
    + creates an empty cache with the given default TTL
    ? entries store an absolute expiry computed at insert time
    # construction
    -> std.sync.new_mutex
  gocache.set
    @ (cache: cache_state, key: string, value: bytes) -> void
    + inserts or replaces an entry using the default TTL
    # mutation
    -> std.time.now_millis
    -> std.sync.with_lock
  gocache.get
    @ (cache: cache_state, key: string) -> optional[bytes]
    + returns the entry when present and not expired
    - returns none when the key is absent
    - returns none when the entry has expired (and removes it)
    # access
    -> std.time.now_millis
    -> std.sync.with_lock
  gocache.delete
    @ (cache: cache_state, key: string) -> void
    + removes an entry if it exists
    # mutation
    -> std.sync.with_lock
  gocache.prune_expired
    @ (cache: cache_state) -> u32
    + removes all expired entries and returns how many were removed
    # maintenance
    -> std.time.now_millis
    -> std.sync.with_lock
