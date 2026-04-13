# Requirement: "an in-memory cache with per-item expiration"

Values are stored with a ttl; reads past the ttl miss.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

cache
  cache.new
    @ () -> cache_state
    + creates an empty cache
    # construction
  cache.set
    @ (c: cache_state, key: string, value: bytes, ttl_millis: i64) -> cache_state
    + stores the value with an absolute expiry of now + ttl_millis
    ? a ttl of zero or negative means the entry is immediately expired
    # write
    -> std.time.now_millis
  cache.get
    @ (c: cache_state, key: string) -> optional[bytes]
    + returns the value when present and not expired
    - returns none when missing
    - returns none when expired
    # read
    -> std.time.now_millis
  cache.evict_expired
    @ (c: cache_state) -> cache_state
    + removes all entries whose expiry is in the past
    # maintenance
    -> std.time.now_millis
