# Requirement: "a cache library with expirable, LFU, LRU, and ARC policies"

A single cache interface with four eviction strategies. Expiry is clock-driven through a std time primitive so tests can control time.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

cache
  cache.new_expirable
    @ (capacity: i32, ttl_ms: i64) -> cache_state
    + creates a cache that evicts entries older than ttl_ms
    # construction
  cache.new_lfu
    @ (capacity: i32) -> cache_state
    + creates a cache that evicts the least-frequently-used entry on overflow
    # construction
  cache.new_lru
    @ (capacity: i32) -> cache_state
    + creates a cache that evicts the least-recently-used entry on overflow
    # construction
  cache.new_arc
    @ (capacity: i32) -> cache_state
    + creates an adaptive cache that balances recency and frequency
    ? maintains T1/T2/B1/B2 lists per the ARC algorithm
    # construction
  cache.set
    @ (state: cache_state, key: string, value: bytes) -> cache_state
    + inserts or updates an entry and evicts per the active policy
    # mutation
    -> std.time.now_millis
  cache.get
    @ (state: cache_state, key: string) -> tuple[optional[bytes], cache_state]
    + returns the value if present and updates bookkeeping (recency/frequency)
    - returns none when key is missing or expired
    # lookup
    -> std.time.now_millis
  cache.delete
    @ (state: cache_state, key: string) -> cache_state
    + removes the entry if present
    # mutation
  cache.len
    @ (state: cache_state) -> i32
    + returns the number of live (non-expired) entries
    # introspection
    -> std.time.now_millis
