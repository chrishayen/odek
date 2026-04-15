# Requirement: "a library for tracing api calls and monitoring api performance, health, and usage metrics"

Record per-call spans in memory and expose aggregated metrics. Time reads go through a thin std utility so tests can substitute a deterministic clock.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.math
    std.math.percentile
      fn (values: list[f64], p: f64) -> f64
      + returns the p-th percentile using linear interpolation (p in [0,1])
      + returns 0 for an empty list
      # statistics

api_monitor
  api_monitor.new
    fn () -> monitor_state
    + returns an empty monitor with zero recorded calls
    # construction
  api_monitor.start_span
    fn (state: monitor_state, route: string) -> tuple[span_id, monitor_state]
    + returns a fresh span id and a state with the pending span registered
    # tracing
    -> std.time.now_millis
  api_monitor.finish_span
    fn (state: monitor_state, id: span_id, status_code: i32) -> monitor_state
    + closes the span and records its duration against the route
    - leaves state unchanged when the id is unknown
    # tracing
    -> std.time.now_millis
  api_monitor.route_stats
    fn (state: monitor_state, route: string) -> route_stats
    + returns count, error count, p50, p95, p99 latency for the route
    + returns zeroed stats when the route has no finished spans
    # metrics
    -> std.math.percentile
  api_monitor.overall_stats
    fn (state: monitor_state) -> overall_stats
    + returns total calls, error rate, and mean latency across all routes
    # metrics
  api_monitor.health
    fn (state: monitor_state, error_rate_threshold: f64) -> health_status
    + returns healthy when the error rate is below the threshold
    - returns degraded otherwise
    # health
  api_monitor.recent_errors
    fn (state: monitor_state, limit: i32) -> list[error_record]
    + returns the most recent spans whose status_code is 500 or above
    + caps the result at limit
    # errors
