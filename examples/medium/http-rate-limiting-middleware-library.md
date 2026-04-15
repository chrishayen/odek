# Requirement: "an HTTP rate limiting middleware library"

Wraps an inner HTTP handler with per-key rate limiting. The key function is pluggable so callers can rate limit by IP, header, or user.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

http_rate_limit
  http_rate_limit.new
    fn (rate_per_sec: f64, burst: i32) -> rate_limit_state
    + creates a limiter with the given refill rate and burst capacity per key
    # construction
  http_rate_limit.allow
    fn (state: rate_limit_state, key: string) -> tuple[bool, rate_limit_state]
    + returns (true, new_state) when the key has tokens available
    - returns (false, unchanged_state) when the key has no tokens
    # rate_limiting
    -> std.time.now_millis
  http_rate_limit.middleware
    fn (state: rate_limit_state, key_fn: fn(http_request) -> string, inner: fn(http_request) -> http_response) -> fn(http_request) -> http_response
    + returns a handler that consults the limiter before calling inner
    + responds with 429 and a Retry-After header when the key is limited
    # middleware
  http_rate_limit.set_response
    fn (state: rate_limit_state, status: i32, body: string) -> rate_limit_state
    + customizes the response returned when a key is limited
    # configuration
