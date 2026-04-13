# Requirement: "a preemptive optimistic cache"

A key-value cache that refreshes entries in the background before they expire so reads rarely miss.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

pocache
  pocache.new
    @ (ttl_ms: i64, refresh_threshold_ms: i64) -> cache_state
    + creates a cache with the given total TTL and early-refresh window
    ? refresh_threshold_ms is the window before expiry during which a hit triggers refresh
    # construction
  pocache.get
    @ (state: cache_state, key: string, loader: loader_fn) -> result[tuple[string, cache_state], string]
    + returns the cached value and schedules a refresh when inside the threshold
    + calls loader on a miss and caches the result
    - returns error when loader fails on a miss
    # read
    -> std.time.now_millis
  pocache.set
    @ (state: cache_state, key: string, value: string) -> cache_state
    + stores value under key with a fresh expiry
    # write
    -> std.time.now_millis
  pocache.invalidate
    @ (state: cache_state, key: string) -> cache_state
    + removes key from the cache
    # eviction
