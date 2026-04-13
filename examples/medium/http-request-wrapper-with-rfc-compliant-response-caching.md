# Requirement: "an HTTP request wrapper with RFC-compliant response caching"

Wraps a request function with a cache that honors Cache-Control, ETag, and Last-Modified semantics.

std
  std.http
    std.http.request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error on transport failure
      # http
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

httpcache
  httpcache.new
    @ (capacity: i32) -> cache_state
    + creates an empty cache with the given maximum entry count
    # construction
  httpcache.parse_cache_control
    @ (header: string) -> cache_control
    + parses a Cache-Control header into a typed struct
    + recognizes max-age, no-store, no-cache, private, public
    # parsing
  httpcache.is_cacheable
    @ (method: string, status: i32, cc: cache_control) -> bool
    + returns true when the response may be stored per RFC 9111
    - returns false for no-store, private, or non-GET methods with no explicit allowance
    # policy
  httpcache.freshness_seconds
    @ (cc: cache_control, date: i64, now: i64) -> i64
    + returns remaining freshness in seconds
    + uses max-age when present, otherwise heuristic based on age
    # policy
    -> std.time.now_seconds
  httpcache.lookup
    @ (state: cache_state, key: string, now: i64) -> optional[http_response]
    + returns a stored response when still fresh
    # retrieval
  httpcache.revalidate
    @ (state: cache_state, key: string) -> result[http_response, string]
    + issues a conditional request using ETag or Last-Modified and updates the entry
    + on 304 refreshes the stored entry's validators and returns the cached body
    - returns error on transport failure
    # revalidation
    -> std.http.request
    -> std.time.now_seconds
  httpcache.request
    @ (state: cache_state, method: string, url: string, headers: map[string, string], body: bytes) -> result[tuple[cache_state, http_response], string]
    + serves from cache when fresh, otherwise fetches, stores, and returns
    - returns error on transport failure
    # pipeline
    -> std.http.request
    -> std.time.now_seconds
