# Requirement: "a query result cache"

A cache keyed on the query string with TTL expiration, sitting in front of a caller-provided query executor.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.crypto
    std.crypto.sha1_hex
      @ (data: bytes) -> string
      + returns the lowercase hex SHA-1 of the input
      # cryptography

query_cache
  query_cache.new
    @ (default_ttl_ms: i64, max_entries: i32) -> cache_state
    + creates a cache with a default TTL and maximum entry count
    # construction
  query_cache.cache_key
    @ (sql: string, params: list[string]) -> string
    + computes a stable key from a query and its bound parameters
    # keying
    -> std.crypto.sha1_hex
  query_cache.get
    @ (cache: cache_state, key: string) -> optional[bytes]
    + returns a non-expired cached value
    - returns none when the entry is missing or expired
    # reads
    -> std.time.now_millis
  query_cache.put
    @ (cache: cache_state, key: string, value: bytes, ttl_ms: i64) -> cache_state
    + stores a value with an explicit TTL, evicting the oldest entry when full
    # writes
    -> std.time.now_millis
  query_cache.invalidate
    @ (cache: cache_state, key: string) -> cache_state
    + removes a specific cached entry
    # invalidation
