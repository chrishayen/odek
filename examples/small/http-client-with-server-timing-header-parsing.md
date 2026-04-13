# Requirement: "an http client that parses the Server-Timing header"

A thin wrapper that performs a request and surfaces Server-Timing metrics.

std
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs a GET request and returns the response with headers
      - returns error on connection failure
      # http

server_timing
  server_timing.fetch_with_metrics
    @ (url: string) -> result[timed_response, string]
    + returns the response plus a list of (name, duration_ms) parsed from Server-Timing
    + returns an empty metrics list when the header is absent
    - returns error when the underlying request fails
    # fetching
    -> std.http.get
  server_timing.parse_header
    @ (header_value: string) -> list[timing_metric]
    + parses comma-separated entries with "name;dur=123" form
    + skips entries that have no "dur=" parameter
    # parsing
