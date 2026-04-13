# Requirement: "an API load testing library"

Drives a sequence of http requests at a configurable concurrency, collecting latency and outcome stats.

std
  std.http
    std.http.request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an http request and returns the response
      - returns error on transport failure or timeout
      # http
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

load_test
  load_test.new_plan
    @ (target_url: string, method: string, total_requests: i32, concurrency: i32) -> result[plan, string]
    + creates a load plan bound to a single endpoint
    - returns error when concurrency is zero or exceeds total_requests
    # plan
  load_test.with_body
    @ (p: plan, body: bytes, headers: map[string, string]) -> plan
    + attaches a request body and headers to the plan
    # plan
  load_test.run
    @ (p: plan) -> result[run_report, string]
    + executes the plan and returns a per-request report
    - returns error when the plan is empty
    # execution
    -> std.http.request
    -> std.time.now_millis
  load_test.summarize
    @ (report: run_report) -> summary
    + computes success count, failure count, mean, p50, p95, p99, and max latency in milliseconds
    # statistics
