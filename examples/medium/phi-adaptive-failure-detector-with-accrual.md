# Requirement: "an adaptive accrual failure detector for distributed systems (phi accrual)"

Tracks heartbeat inter-arrival times per node and produces a phi suspicion value the caller can threshold. Time is injected so tests are deterministic.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the square root
      # math
    std.math.exp
      fn (x: f64) -> f64
      + returns e^x
      # math
    std.math.log10
      fn (x: f64) -> f64
      + returns base-10 logarithm
      # math

phi
  phi.detector_new
    fn (window_size: i32, min_std_dev_ms: f64, initial_interval_ms: f64) -> detector_state
    + creates a detector with a sliding window of the given size
    ? min_std_dev_ms prevents division-by-zero and sets a floor on variability
    # construction
  phi.heartbeat
    fn (state: detector_state, node_id: string) -> detector_state
    + records that a heartbeat from node_id arrived now
    + updates the node's sliding window of inter-arrival times
    # ingest
    -> std.time.now_millis
  phi.value
    fn (state: detector_state, node_id: string) -> f64
    + returns the current phi suspicion value for the node
    ? zero means fully trusted; values rise with elapsed time since last heartbeat
    + returns zero for nodes with no prior heartbeats
    # suspicion
    -> std.time.now_millis
    -> std.math.exp
    -> std.math.log10
    -> std.math.sqrt
  phi.is_available
    fn (state: detector_state, node_id: string, threshold: f64) -> bool
    + returns true when the current phi value is below the threshold
    # decision
    -> phi.value
  phi.forget
    fn (state: detector_state, node_id: string) -> detector_state
    + drops all state for the given node
    # lifecycle
