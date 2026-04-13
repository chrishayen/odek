# Requirement: "a map where entries expire after a configurable lifetime"

Per-entry TTLs with lazy expiration on read. Time reads go through a std primitive so tests can substitute a clock.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

timedmap
  timedmap.new
    @ () -> timedmap_state
    + creates an empty map
    # construction
  timedmap.set
    @ (state: timedmap_state, key: string, value: bytes, ttl_ms: i64) -> timedmap_state
    + stores the value with an expiration ttl_ms milliseconds in the future
    + overwrites any existing entry for key
    # mutation
    -> std.time.now_millis
  timedmap.get
    @ (state: timedmap_state, key: string) -> tuple[optional[bytes], timedmap_state]
    + returns the value when the key is present and not expired
    + returns none and drops the entry when it has expired
    - returns none when the key was never set
    # access
    -> std.time.now_millis
  timedmap.remove
    @ (state: timedmap_state, key: string) -> timedmap_state
    + removes the entry for key if present
    # mutation
  timedmap.sweep
    @ (state: timedmap_state) -> timedmap_state
    + drops every entry whose expiration has passed
    # maintenance
    -> std.time.now_millis
