# Requirement: "an in-memory key-value cache with automatic timeout-based invalidation"

Each entry has an idle timeout that resets on access.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

cache
  cache.new
    fn (default_timeout_millis: i64) -> cache_state
    + creates an empty cache with the given default idle timeout
    # construction
  cache.put
    fn (c: cache_state, key: string, value: bytes) -> cache_state
    + stores the value and sets its last-access timestamp to now
    # write
    -> std.time.now_millis
  cache.lookup
    fn (c: cache_state, key: string) -> tuple[optional[bytes], cache_state]
    + returns the value and a state with a refreshed last-access timestamp
    - returns none when missing or idle-expired
    # read
    -> std.time.now_millis
  cache.delete
    fn (c: cache_state, key: string) -> cache_state
    + removes the key if present
    # write
  cache.prune_idle
    fn (c: cache_state) -> cache_state
    + drops entries whose idle time exceeds the default timeout
    # maintenance
    -> std.time.now_millis
