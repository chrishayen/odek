# Requirement: "a generic thread-safe in-memory cache"

A keyed cache holding opaque byte payloads. The cache is internally synchronized so concurrent callers can safely read and write.

std
  std.concurrency
    std.concurrency.mutex_new
      @ () -> mutex_state
      + returns a fresh unlocked mutex
      # concurrency

cache
  cache.new
    @ (capacity: i32) -> cache_state
    + returns an empty cache with the given maximum entry count
    ? capacity <= 0 means unbounded
    # construction
    -> std.concurrency.mutex_new
  cache.get
    @ (cache: cache_state, key: string) -> optional[bytes]
    + returns the stored value when the key is present
    # reads
  cache.set
    @ (cache: cache_state, key: string, value: bytes) -> cache_state
    + stores the value, evicting the least recently used entry when at capacity
    # writes
  cache.delete
    @ (cache: cache_state, key: string) -> cache_state
    + removes the entry for the key; no-op when absent
    # writes
  cache.size
    @ (cache: cache_state) -> i32
    + returns the current number of stored entries
    # introspection
