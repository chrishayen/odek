# Requirement: "a live runtime-statistics visualization library"

Collects sampled runtime metrics and exposes them as a time series that a UI layer can render.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.json
    std.json.encode_object
      fn (obj: map[string, f64]) -> string
      + encodes a string-to-float map as a JSON object
      # serialization

runtime_viz
  runtime_viz.new
    fn (max_samples: i32) -> viz_state
    + creates a ring buffer that retains the most recent N samples per metric
    ? older samples are discarded once the ring is full
    # construction
  runtime_viz.record
    fn (state: viz_state, metric: string, value: f64) -> viz_state
    + appends a timestamped sample for the named metric
    # sampling
    -> std.time.now_millis
  runtime_viz.record_many
    fn (state: viz_state, metrics: map[string, f64]) -> viz_state
    + appends one timestamped sample for each metric in a single tick
    # sampling
    -> std.time.now_millis
  runtime_viz.series
    fn (state: viz_state, metric: string) -> list[tuple[i64, f64]]
    + returns the retained (timestamp, value) pairs for the named metric in chronological order
    - returns an empty list when the metric has never been recorded
    # query
  runtime_viz.snapshot_json
    fn (state: viz_state) -> string
    + encodes the latest value of every recorded metric as a JSON object
    # export
    -> std.json.encode_object
