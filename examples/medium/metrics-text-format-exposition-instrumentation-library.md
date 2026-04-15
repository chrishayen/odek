# Requirement: "a metrics instrumentation library with text-format exposition"

Counter, gauge, and histogram primitives plus an exposition format suitable for scraping endpoints.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current time in milliseconds
      # time

metrics
  metrics.new_registry
    fn () -> registry_state
    + creates an empty metric registry
    # construction
  metrics.counter
    fn (registry: registry_state, name: string, labels: map[string,string]) -> registry_state
    + increments a counter identified by name and label set
    ? missing counters are created with value 0 on first touch
    # counter
  metrics.add
    fn (registry: registry_state, name: string, labels: map[string,string], delta: f64) -> registry_state
    + adds delta to the named counter
    - clamps delta to non-negative for counters
    # counter
  metrics.set_gauge
    fn (registry: registry_state, name: string, labels: map[string,string], value: f64) -> registry_state
    + sets the gauge to an absolute value
    # gauge
  metrics.observe
    fn (registry: registry_state, name: string, labels: map[string,string], value: f64) -> registry_state
    + records an observation into a histogram with default buckets
    # histogram
  metrics.expose_text
    fn (registry: registry_state) -> string
    + renders all metrics in scrape-compatible text format
    + one metric line per series with label set
    # exposition
  metrics.time_block
    fn (registry: registry_state, name: string, labels: map[string,string], start_ms: i64) -> registry_state
    + records elapsed milliseconds since start_ms into a histogram
    # timing
    -> std.time.now_millis
