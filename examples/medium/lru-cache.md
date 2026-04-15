# Requirement: "an LRU cache"

A classic data structure — three operations on an opaque state value. Tests cover the eviction policy directly on `put`.

std: (all units exist)

lru_cache
  lru_cache.new
    fn (capacity: i32) -> lru_cache_state
    + creates an empty cache with the given capacity
    ? capacity must be >= 1; validating that is the caller's job
    # construction
  lru_cache.get
    fn (state: lru_cache_state, key: string) -> tuple[optional[string], lru_cache_state]
    + returns (some(value), new_state) when the key is present and marks it as recently used
    + returns (none, unchanged_state) when the key is absent
    # cache_access
  lru_cache.put
    fn (state: lru_cache_state, key: string, value: string) -> lru_cache_state
    + inserts the value and marks it as most recently used
    + evicts the least recently used entry when at capacity
    + updating an existing key refreshes its position and does not evict
    # cache_access
