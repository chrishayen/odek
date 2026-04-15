# Requirement: "a type-safe metrics builder wrapper"

Wraps a metrics registry with strongly-typed counter, gauge, and histogram handles so label sets are validated at construction time instead of call time.

std: (all units exist)

metrics_builder
  metrics_builder.new_registry
    fn () -> registry
    + returns an empty registry
    # construction
  metrics_builder.build_counter
    fn (reg: registry, name: string, help: string, labels: list[string]) -> result[tuple[counter_handle, registry], string]
    + registers a counter with the given label schema
    - returns error when a metric with the name already exists
    - returns error when labels contain duplicates
    # counter_registration
  metrics_builder.build_gauge
    fn (reg: registry, name: string, help: string, labels: list[string]) -> result[tuple[gauge_handle, registry], string]
    + registers a gauge with the given label schema
    - returns error on a name collision
    # gauge_registration
  metrics_builder.build_histogram
    fn (reg: registry, name: string, help: string, labels: list[string], buckets: list[f64]) -> result[tuple[histogram_handle, registry], string]
    + registers a histogram with the given buckets in ascending order
    - returns error when buckets are not strictly ascending
    # histogram_registration
  metrics_builder.counter_inc
    fn (reg: registry, h: counter_handle, values: list[string], delta: f64) -> result[registry, string]
    + increments the counter series keyed by values
    - returns error when values length does not match the label schema
    # counter_update
  metrics_builder.gauge_set
    fn (reg: registry, h: gauge_handle, values: list[string], x: f64) -> result[registry, string]
    + sets the gauge series keyed by values to x
    - returns error on a label-arity mismatch
    # gauge_update
  metrics_builder.histogram_observe
    fn (reg: registry, h: histogram_handle, values: list[string], x: f64) -> result[registry, string]
    + records an observation into the matching bucket
    - returns error on a label-arity mismatch
    # histogram_update
