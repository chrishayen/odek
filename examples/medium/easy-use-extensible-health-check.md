# Requirement: "an extensible health check library"

Registers named checks, runs them, and aggregates pass/fail results into an overall status.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

health
  health.new
    @ () -> health_registry
    + creates an empty health check registry
    # construction
  health.register
    @ (registry: health_registry, name: string, check: fn() -> result[void, string]) -> result[health_registry, string]
    + stores a named check in the registry
    - returns error when a check with the same name is already registered
    # registration
  health.run_all
    @ (registry: health_registry) -> health_report
    + runs every registered check and collects individual outcomes
    + records duration in milliseconds for each check
    # execution
    -> std.time.now_millis
  health.overall_status
    @ (report: health_report) -> string
    + returns "healthy" when every check passed
    + returns "degraded" when at least one non-critical check failed
    - returns "unhealthy" when any critical check failed
    # aggregation
  health.mark_critical
    @ (registry: health_registry, name: string) -> result[health_registry, string]
    + marks a registered check as critical so its failure forces unhealthy
    - returns error when name is not registered
    # criticality
