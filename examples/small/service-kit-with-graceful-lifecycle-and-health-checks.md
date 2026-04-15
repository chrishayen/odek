# Requirement: "a small services toolkit for graceful lifecycle and health checks"

The original input is a grab-bag of utilities; the coherent core is a service lifecycle with health reporting. Keep it to that.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

service_kit
  service_kit.new
    fn (name: string) -> service_state
    + creates a service in "starting" status with no registered checks
    # construction
    -> std.time.now_seconds
  service_kit.register_check
    fn (state: service_state, name: string, check_id: string) -> service_state
    + adds a named health check bound to check_id
    # health
  service_kit.mark_ready
    fn (state: service_state) -> service_state
    + transitions the service from "starting" to "ready"
    # lifecycle
  service_kit.begin_shutdown
    fn (state: service_state) -> service_state
    + transitions the service to "draining" and refuses further ready transitions
    # lifecycle
  service_kit.report
    fn (state: service_state, results: map[string, bool]) -> health_snapshot
    + returns a snapshot containing service status, uptime seconds, and per-check pass/fail from results
    # reporting
    -> std.time.now_seconds
