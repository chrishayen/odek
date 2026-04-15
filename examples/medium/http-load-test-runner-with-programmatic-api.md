# Requirement: "a library for running HTTP load tests with a programmatic API"

Define a scenario (requests, concurrency, duration or request count), run it, and aggregate latency and status statistics.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs a single HTTP request and returns the response
      - returns error on connection failure
      # http
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns current monotonic time in nanoseconds
      # time
  std.concurrency
    std.concurrency.run_workers
      fn (worker_count: i32, task: fn(i32) -> void) -> void
      + runs task concurrently across worker_count workers and waits for all to finish
      # concurrency

loadtest
  loadtest.new_scenario
    fn (name: string) -> scenario
    + creates an empty scenario with the given name
    # construction
  loadtest.add_step
    fn (s: scenario, method: string, url: string, headers: map[string, string], body: bytes) -> scenario
    + appends a request step; workers iterate the step list in order
    # scenario
  loadtest.configure
    fn (s: scenario, concurrency: i32, duration_ms: i64, max_requests: i32) -> scenario
    + sets the worker count and stop condition; a zero duration means "until max_requests is reached"
    # scenario
  loadtest.run
    fn (s: scenario) -> run_report
    + executes the scenario, measuring latency per request and counting outcomes per status class
    # execution
    -> std.concurrency.run_workers
    -> std.http.request
    -> std.time.now_nanos
  loadtest.summarize
    fn (report: run_report) -> summary
    + computes mean, p50, p95, p99 latency, requests per second, and error rate
    # reporting
