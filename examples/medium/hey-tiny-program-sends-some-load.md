# Requirement: "a load generator that sends concurrent HTTP requests to a web endpoint and reports latency statistics"

std
  std.net
    std.net.http_request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs a single HTTP request and returns the response
      - returns error on transport failure
      # networking
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns current monotonic time in nanoseconds
      # time
  std.math
    std.math.sort_f64
      @ (values: list[f64]) -> list[f64]
      + returns values sorted in ascending order
      # math

load_gen
  load_gen.new_plan
    @ (method: string, url: string, total_requests: i32, concurrency: i32) -> load_plan
    + builds a plan describing the load to generate
    ? concurrency is clamped to total_requests
    # planning
  load_gen.run
    @ (plan: load_plan) -> load_report
    + issues the requests with the configured concurrency and records outcomes
    + each outcome captures latency in nanoseconds and status code
    # execution
    -> std.net.http_request
    -> std.time.now_nanos
  load_gen.latency_percentile
    @ (report: load_report, percentile: f64) -> f64
    + returns the given percentile of latency in milliseconds
    ? percentile is a value in [0, 100]
    # statistics
    -> std.math.sort_f64
  load_gen.summary
    @ (report: load_report) -> load_summary
    + returns total requests, error count, mean latency, min, max, and p50/p95/p99
    # statistics
    -> load_gen.latency_percentile
  load_gen.status_histogram
    @ (report: load_report) -> map[i32, i32]
    + returns a count of responses grouped by HTTP status code
    # statistics
