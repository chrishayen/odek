# Requirement: "a concurrent health-check http handler for services"

Register named checks, run them concurrently, and expose the result as an http response body plus status.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.concurrent
    std.concurrent.run_parallel
      fn (tasks: list[task_handle]) -> list[string]
      + runs tasks in parallel and returns their results in input order
      # concurrency

health
  health.new_registry
    fn () -> registry_state
    + returns an empty registry
    # construction
  health.register_liveness
    fn (r: registry_state, name: string, check: check_fn) -> registry_state
    + adds a liveness check; duplicate names replace the previous
    # registration
  health.register_readiness
    fn (r: registry_state, name: string, check: check_fn) -> registry_state
    + adds a readiness check; duplicate names replace the previous
    # registration
  health.run_all
    fn (r: registry_state, timeout_ms: i64) -> health_report
    + runs every check concurrently and returns a report with per-check status
    + marks a check failed when its result does not arrive within timeout_ms
    # execution
    -> std.concurrent.run_parallel
    -> std.time.now_millis
  health.report_status_code
    fn (report: health_report) -> i32
    + returns 200 when all checks passed, 503 otherwise
    # http
  health.report_body
    fn (report: health_report) -> string
    + returns a JSON body with overall status and per-check detail
    # http
    -> std.json.encode_object
  health.handle_request
    fn (r: registry_state, path: string, timeout_ms: i64) -> tuple[i32, string]
    + dispatches "/livez" and "/readyz" to the corresponding check set
    - returns (404, "") for any other path
    # http
