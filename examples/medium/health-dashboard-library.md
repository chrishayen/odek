# Requirement: "a service health dashboard library"

Runs health checks against registered endpoints on a schedule and exposes a rollup of results. HTTP and time are std primitives; the project holds the check registry and status logic.

std
  std.http
    std.http.get
      fn (url: string, timeout_ms: i32) -> result[http_response, string]
      + returns status code, headers, and body
      - returns error on connection failure or timeout
      # http_client
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

health
  health.register
    fn (dash: dashboard_state, name: string, url: string, interval_ms: i32) -> dashboard_state
    + adds a new endpoint to the dashboard with its check interval
    # registration
  health.check_once
    fn (endpoint: endpoint_config) -> check_result
    + returns ok when the endpoint responds with 2xx within its timeout
    - returns failing when the response status is not 2xx
    - returns failing when the request errors out
    # probe
    -> std.http.get
    -> std.time.now_millis
  health.tick
    fn (dash: dashboard_state) -> dashboard_state
    + runs any endpoint whose last-check time is older than its interval
    + stores the latest result on each endpoint
    # scheduling
    -> std.time.now_millis
  health.status
    fn (dash: dashboard_state) -> rollup
    + returns a summary with counts of healthy, degraded, and failing endpoints
    # reporting
  health.endpoint_history
    fn (dash: dashboard_state, name: string) -> result[list[check_result], string]
    + returns the recent check results for a named endpoint
    - returns error when the endpoint is unknown
    # history
