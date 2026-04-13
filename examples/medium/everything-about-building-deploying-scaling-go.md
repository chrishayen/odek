# Requirement: "a library of utilities supporting building, deploying, and scaling applications in production"

Production support helpers: a health registry, graceful shutdown coordination, a basic metrics registry, and a simple rolling-window counter.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

prodkit
  prodkit.health_registry
    @ () -> health_state
    + creates an empty health registry
    # construction
  prodkit.register_check
    @ (state: health_state, name: string, check_tag: string) -> health_state
    + registers a named health check referenced by tag
    # health
  prodkit.run_checks
    @ (state: health_state) -> list[check_result]
    + runs every registered check and returns the results
    + an overall status of healthy when all checks pass
    - overall status of unhealthy when any check fails
    # health
  prodkit.shutdown_state
    @ () -> shutdown_state
    + creates a shutdown coordinator with no hooks registered
    # construction
  prodkit.on_shutdown
    @ (state: shutdown_state, hook_tag: string) -> shutdown_state
    + registers a hook to run during graceful shutdown
    # shutdown
  prodkit.begin_shutdown
    @ (state: shutdown_state, deadline_millis: i64) -> shutdown_state
    + marks the system as draining with an absolute deadline
    -> std.time.now_millis
    # shutdown
  prodkit.is_draining
    @ (state: shutdown_state) -> bool
    + returns true when shutdown has been initiated
    # shutdown
  prodkit.metrics_registry
    @ () -> metrics_state
    + creates an empty metrics registry
    # construction
  prodkit.counter_inc
    @ (state: metrics_state, name: string, amount: i64) -> metrics_state
    + increases a named counter by the given amount
    # metrics
  prodkit.gauge_set
    @ (state: metrics_state, name: string, value: f64) -> metrics_state
    + sets a named gauge to the given value
    # metrics
  prodkit.rolling_window
    @ (size_seconds: i32) -> rolling_state
    + creates a rolling-window counter of the given span
    # construction
  prodkit.rolling_record
    @ (state: rolling_state, value: i64) -> rolling_state
    + records a value at the current time, expiring entries outside the window
    -> std.time.now_millis
    # metrics
  prodkit.rolling_sum
    @ (state: rolling_state) -> i64
    + returns the sum of values currently inside the window
    # metrics
