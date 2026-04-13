# Requirement: "an HTTP load testing library"

The library schedules requests at a given rate, records latencies, and summarizes the results. Making the requests themselves is delegated to a std HTTP primitive.

std
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs an HTTP GET and returns status, headers, and body
      - returns error on network failure
      # http
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time

loadtest
  loadtest.new_run
    @ (url: string, total_requests: i32) -> run_state
    + creates a run targeting url with the given request count and an empty sample buffer
    # construction
  loadtest.issue_one
    @ (state: run_state) -> run_state
    + performs one request and appends its latency and status to the run
    # execution
    -> std.http.get
    -> std.time.now_nanos
  loadtest.percentile
    @ (state: run_state, p: f64) -> i64
    + returns the pth-percentile latency in nanoseconds
    - returns 0 when no samples have been recorded
    # stats
  loadtest.error_rate
    @ (state: run_state) -> f64
    + returns the fraction of requests with status >= 400 or network errors
    # stats
  loadtest.summary
    @ (state: run_state) -> load_summary
    + returns counts, error rate, and p50/p95/p99 latencies
    # reporting
