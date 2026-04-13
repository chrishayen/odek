# Requirement: "a read/write splitting pool for a primary database with multiple read replicas"

One primary for writes, many replicas for reads, with round-robin selection and a simple health gate.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

rwdb
  rwdb.new
    @ (primary_dsn: string, replica_dsns: list[string]) -> rwdb_pool
    + creates a pool with one primary and zero or more replicas
    ? replicas list may be empty; in that case reads fall back to the primary
    # construction
  rwdb.pick_writer
    @ (pool: rwdb_pool) -> string
    + returns the primary DSN
    # routing
  rwdb.pick_reader
    @ (pool: rwdb_pool) -> tuple[string, rwdb_pool]
    + returns the next healthy replica DSN using round-robin and advances the cursor
    + falls back to the primary when no replica is healthy
    # routing
  rwdb.mark_unhealthy
    @ (pool: rwdb_pool, dsn: string) -> rwdb_pool
    + records a replica as unhealthy and excludes it from selection
    # health
    -> std.time.now_millis
  rwdb.recover
    @ (pool: rwdb_pool, cooldown_ms: i64) -> rwdb_pool
    + re-enables replicas whose cooldown has elapsed since being marked unhealthy
    # health
    -> std.time.now_millis
