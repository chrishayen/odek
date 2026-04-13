# Requirement: "a reverse-proxy cache tuned for time-series queries"

Caches time-ranged query responses and serves overlapping requests from the cache, delta-fetching only the missing range.

std
  std.http
    std.http.forward
      @ (url: string, body: bytes) -> result[http_response, string]
      + forwards a request to the upstream and returns the response
      - returns error on connection failure
      # http
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

ts_cache
  ts_cache.new
    @ (max_entries: i32) -> cache_state
    + returns an empty cache with the given capacity
    # construction
  ts_cache.lookup
    @ (state: cache_state, key: string, start_ms: i64, end_ms: i64) -> cache_hit
    + returns the full cached series when the requested range is fully covered
    + returns the covered sub-range and the missing sub-range when partially covered
    - returns a miss indicator when the key is absent
    # lookup
    -> std.time.now_millis
  ts_cache.store
    @ (state: cache_state, key: string, start_ms: i64, end_ms: i64, series: list[sample]) -> cache_state
    + stores the range as a single entry, merging with adjacent ranges under the same key
    + evicts the least-recently-used entry when over capacity
    # storage
  ts_cache.handle_request
    @ (state: cache_state, upstream_url: string, key: string, start_ms: i64, end_ms: i64) -> result[tuple[list[sample], cache_state], string]
    + serves from cache when covered, otherwise forwards the missing range and merges the result
    - returns error when the upstream fetch fails
    # dispatch
    -> std.http.forward
  ts_cache.merge_ranges
    @ (existing: list[sample], delta: list[sample]) -> list[sample]
    + returns a single sorted series with duplicates by timestamp removed
    # merging
