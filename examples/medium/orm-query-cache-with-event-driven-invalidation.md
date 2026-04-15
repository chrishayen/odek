# Requirement: "an ORM query cache with event-driven invalidation"

A cache keyed by query fingerprint. When the caller reports a mutation on a table, every cached entry whose dependencies touch that table is dropped.

std
  std.hash
    std.hash.sha1_hex
      fn (data: bytes) -> string
      + returns the SHA-1 digest as a lowercase hex string
      # hashing
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

query_cache
  query_cache.new
    fn (default_ttl_millis: i64) -> cache_state
    + creates an empty cache with the given default TTL
    # construction
  query_cache.fingerprint
    fn (sql: string, params: list[string]) -> string
    + returns a stable fingerprint for a SQL statement and its bound parameters
    + fingerprints are equal for identical sql and params
    # keying
    -> std.hash.sha1_hex
  query_cache.put
    fn (state: cache_state, fp: string, tables: list[string], value: bytes, now_millis: i64) -> cache_state
    + stores value under fp and records the tables the query depends on
    + overwrites any prior entry for the same fp
    # storage
  query_cache.get
    fn (state: cache_state, fp: string, now_millis: i64) -> optional[bytes]
    + returns the stored value when present and not expired
    - returns none when the entry has exceeded the default TTL
    # lookup
    -> std.time.now_millis
  query_cache.invalidate_table
    fn (state: cache_state, table: string) -> cache_state
    + removes every entry whose dependency list contains the table
    + returns unchanged state when no entries depend on the table
    # invalidation
  query_cache.size
    fn (state: cache_state) -> i32
    + returns the number of live entries
    # inspection
