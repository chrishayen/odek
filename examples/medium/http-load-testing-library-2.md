# Requirement: "an http load testing library"

Drives requests at a target rate and reports latency and error statistics.

std
  std.http
    std.http.request
      fn (method: string, url: string, body: bytes) -> result[http_response, string]
      + performs a request and returns the response
      - returns error on connection or timeout failure
      # http
  std.time
    std.time.now_micros
      fn () -> i64
      + returns current time as a unix microsecond count
      # time

load_test
  load_test.new_plan
    fn (target_url: string, method: string, rate_per_sec: f64, duration_sec: i32) -> attack_plan
    + creates a plan describing the load profile
    ? request body is empty; callers extend the plan if needed
    # construction
  load_test.run
    fn (plan: attack_plan) -> attack_report
    + issues requests on the schedule dictated by rate and duration
    + records a latency sample and status for every issued request
    # execution
    -> std.http.request
    -> std.time.now_micros
  load_test.summarize
    fn (report: attack_report) -> load_summary
    + computes count, error_count, mean, p50, p95, p99 latency in microseconds
    + returns a zeroed summary when the report has no samples
    # analysis
