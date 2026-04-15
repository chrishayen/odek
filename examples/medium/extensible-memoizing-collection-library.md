# Requirement: "an extensible memoizing collection library"

A core cache store plus pluggable eviction policies. The project exposes a handful of caches sharing one lookup primitive.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.hash
    std.hash.hash_string
      fn (s: string) -> u64
      + returns a stable 64-bit hash of the input
      # hashing

cache
  cache.new_lru
    fn (max_entries: i32) -> cache_state
    + creates an empty cache that evicts the least recently used entry when full
    ? entries are keyed and valued as strings; callers serialize their own types
    # construction
  cache.new_lfu
    fn (max_entries: i32) -> cache_state
    + creates an empty cache that evicts the least frequently used entry when full
    # construction
  cache.new_ttl
    fn (max_entries: i32, ttl_millis: i64) -> cache_state
    + creates an empty cache where entries expire after the given time-to-live
    # construction
    -> std.time.now_millis
  cache.get
    fn (state: cache_state, key: string) -> tuple[optional[string], cache_state]
    + returns (Some(value), new_state) on hit and updates recency or frequency metadata
    - returns (None, new_state) on miss or when the entry has expired
    # lookup
    -> std.time.now_millis
  cache.put
    fn (state: cache_state, key: string, value: string) -> cache_state
    + inserts or replaces an entry and evicts one entry if at capacity
    # insertion
    -> std.time.now_millis
  cache.memoize
    fn (state: cache_state, key: string, compute: fn(string) -> string) -> tuple[string, cache_state]
    + returns the cached value when present; otherwise calls compute, stores, and returns the result
    # memoization
