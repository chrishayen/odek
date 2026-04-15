# Requirement: "an in-memory cache with proactive TTL expiration"

A typed key/value cache with per-entry time-to-live. Expired entries are proactively evicted via a background expiry wheel rather than lazily on access.

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
    fn (capacity: i64) -> cache_state
    + creates a cache with the given entry capacity
    # construction
  cache.set
    fn (state: cache_state, key: string, value: bytes, ttl_ms: i64) -> cache_state
    + inserts or updates an entry with the given TTL
    # insertion
    -> std.hash.hash64
    -> std.time.now_millis
  cache.get
    fn (state: cache_state, key: string) -> optional[bytes]
    + returns the value when present and not expired
    - returns none when the key is missing or has expired
    # lookup
    -> std.hash.hash64
    -> std.time.now_millis
  cache.delete
    fn (state: cache_state, key: string) -> cache_state
    + removes the entry when present
    # removal
  cache.tick
    fn (state: cache_state) -> tuple[i64, cache_state]
    + advances the expiry wheel and returns the number of entries purged
    # expiry
    -> std.time.now_millis
  cache.len
    fn (state: cache_state) -> i64
    + returns the number of live entries
    # observability
