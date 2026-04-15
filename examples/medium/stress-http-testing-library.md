# Requirement: "an HTTP stress testing library"

Fires configurable numbers of concurrent HTTP requests and reports latency and error statistics.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs the request and returns status, headers, and body
      - returns error on connection failure or timeout
      # http
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns current monotonic time in nanoseconds
      # time
  std.concurrency
    std.concurrency.run_parallel
      fn (worker_count: i32, task_count: i64, task: parallel_task) -> list[task_result]
      + runs task_count tasks across worker_count workers, returning results in completion order
      # concurrency

stress
  stress.new_plan
    fn (url: string, method: string) -> stress_plan
    + creates a plan targeting a URL with the given method
    # construction
  stress.with_requests
    fn (plan: stress_plan, total: i64, concurrency: i32) -> stress_plan
    + sets the total request count and worker concurrency
    # config
  stress.with_header
    fn (plan: stress_plan, name: string, value: string) -> stress_plan
    + adds a request header to every call
    # config
  stress.with_body
    fn (plan: stress_plan, body: bytes) -> stress_plan
    + sets the request body for every call
    # config
  stress.run
    fn (plan: stress_plan) -> stress_report
    + executes the plan and collects per-request timings and status codes
    # execution
    -> std.http.request
    -> std.time.now_nanos
    -> std.concurrency.run_parallel
  stress.summarize
    fn (report: stress_report) -> stress_summary
    + computes p50, p95, p99, mean latency, total duration, and error count
    + returns zero values when the report has no samples
    # statistics
