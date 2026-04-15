# Requirement: "middleware that collects runtime statistics about a web application"

Tracks request counts, status code breakdowns, and latency per route. The caller decides how to expose the snapshot.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

stats
  stats.new
    fn () -> stats_state
    + creates an empty stats collector
    # construction
  stats.record
    fn (state: stats_state, route: string, status: i32, latency_ms: i64) -> stats_state
    + increments the total counter, bumps the per-status tally, and adds to the route's latency bucket
    # recording
  stats.wrap
    fn (state: stats_state, handler: fn(request) -> response) -> fn(request) -> response
    + returns a handler that times each call and records the outcome
    # middleware
    -> std.time.now_millis
  stats.snapshot
    fn (state: stats_state) -> stats_snapshot
    + returns a copy containing totals, per-status counts, and per-route averages
    # reporting
