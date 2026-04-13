# Requirement: "an HTTP cache facade backed by a pluggable key-value store"

Exposes get, put, and delete as request handlers against a caller-supplied storage backend, with TTL support and optional value compression.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.compress
    std.compress.gzip
      @ (data: bytes) -> bytes
      + returns a gzip-compressed copy of the input
      # compression
    std.compress.gunzip
      @ (data: bytes) -> result[bytes, string]
      + returns the decompressed payload
      - returns error on a malformed stream
      # compression

cache_api
  cache_api.new
    @ (store: kv_store) -> cache_state
    + creates a cache facade over the store
    # construction
  cache_api.handle_get
    @ (state: cache_state, key: string) -> http_response
    + returns 200 with the stored value when present and unexpired
    - returns 404 when the key is absent or expired
    # handler
    -> std.time.now_seconds
    -> std.compress.gunzip
  cache_api.handle_put
    @ (state: cache_state, key: string, body: bytes, ttl_seconds: i64) -> http_response
    + stores the value with the given ttl and returns 204
    - returns 413 when the body exceeds the configured maximum size
    # handler
    -> std.time.now_seconds
    -> std.compress.gzip
  cache_api.handle_delete
    @ (state: cache_state, key: string) -> http_response
    + removes the key and returns 204
    # handler
  cache_api.sweep_expired
    @ (state: cache_state) -> i32
    + removes every expired entry and returns how many were dropped
    # maintenance
    -> std.time.now_seconds
