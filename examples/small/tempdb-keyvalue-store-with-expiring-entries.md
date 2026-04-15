# Requirement: "a key-value store for temporary items with expiring entries"

An in-memory store where each entry carries a deadline. Reads transparently evict expired items.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

tempdb
  tempdb.new
    fn () -> tempdb_state
    + creates an empty store
    # construction
  tempdb.set
    fn (state: tempdb_state, key: string, value: string, ttl_seconds: i64) -> tempdb_state
    + stores the value with an expiry computed from now + ttl
    # write
    -> std.time.now_seconds
  tempdb.get
    fn (state: tempdb_state, key: string) -> optional[string]
    + returns the value when the key exists and has not expired
    - returns none when the key is missing
    - returns none when the entry has expired
    # read
    -> std.time.now_seconds
  tempdb.delete
    fn (state: tempdb_state, key: string) -> tempdb_state
    + removes the key if present
    # write
  tempdb.purge
    fn (state: tempdb_state) -> tempdb_state
    + removes all expired entries
    # maintenance
    -> std.time.now_seconds
