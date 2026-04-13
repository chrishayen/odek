# Requirement: "a benchmarking library for HTTP, TLS, and DNS workloads"

Drives configurable request loads against pluggable targets and aggregates timings.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns the current monotonic time in nanoseconds
      # time
  std.http
    std.http.request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + sends an HTTP request and returns the response
      - returns error on connection, TLS, or protocol failure
      # networking
  std.net
    std.net.tls_handshake
      @ (host: string, port: i32) -> result[tls_info, string]
      + performs a TLS handshake and returns the negotiated parameters
      - returns error on handshake failure
      # networking
    std.net.dns_resolve
      @ (name: string) -> result[list[string], string]
      + returns the list of addresses for name
      - returns error on NXDOMAIN or timeout
      # networking

benchmark
  benchmark.new_config
    @ (concurrency: i32, total_requests: i64, duration_seconds: i32) -> bench_config
    + returns a configuration with the given shape
    ? total_requests of 0 means "run until duration expires"
    # configuration
  benchmark.run_http
    @ (cfg: bench_config, target: http_target) -> bench_report
    + issues HTTP requests according to cfg and records per-request latencies
    # execution
    -> std.http.request
    -> std.time.now_nanos
  benchmark.run_tls
    @ (cfg: bench_config, host: string, port: i32) -> bench_report
    + repeatedly performs TLS handshakes and records their durations
    # execution
    -> std.net.tls_handshake
    -> std.time.now_nanos
  benchmark.run_dns
    @ (cfg: bench_config, name: string) -> bench_report
    + repeatedly resolves name and records per-query durations
    # execution
    -> std.net.dns_resolve
    -> std.time.now_nanos
  benchmark.summarize
    @ (report: bench_report) -> bench_summary
    + returns totals, mean, p50, p95, p99, and error count
    # statistics
  benchmark.format_summary
    @ (summary: bench_summary) -> string
    + returns a human-readable multi-line summary
    # reporting
