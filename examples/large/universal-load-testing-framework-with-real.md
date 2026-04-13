# Requirement: "a load testing framework that runs a user-defined scenario at a target rate, records per-request metrics, and exposes live statistics"

Scenarios are functions producing a single request attempt. The runner drives concurrency and aggregates latency, error counts, and throughput into a live snapshot.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time
    std.time.sleep_nanos
      @ (ns: i64) -> void
      + suspends the caller
      # time
  std.sync
    std.sync.spawn
      @ (fn: fn() -> void) -> void
      + runs fn on a background worker
      # concurrency
    std.sync.atomic_add_i64
      @ (cell: atomic_i64, delta: i64) -> i64
      + atomically adds and returns the new value
      # concurrency

load_test
  load_test.new
    @ (target_rps: f64, concurrency: i32, duration_ms: i64) -> run_state
    + returns a run with the requested rate, worker count, and duration
    # construction
  load_test.set_scenario
    @ (state: run_state, scenario: fn() -> result[void,string]) -> run_state
    + installs the per-request function the workers invoke
    # configuration
  load_test.record_sample
    @ (state: run_state, latency_ns: i64, ok: bool) -> run_state
    + updates histogram buckets, success/failure counts, and running total
    # metrics
  load_test.snapshot
    @ (state: run_state) -> metrics_snapshot
    + returns counts, percentiles (p50/p95/p99), and current throughput
    # metrics
  load_test.rate_gate
    @ (state: run_state, sent: i64, start_ns: i64) -> i64
    + returns how long a worker should sleep to keep the global rate near target
    ? uses a simple pacing computation: expected_time - elapsed
    # pacing
    -> std.time.now_nanos
  load_test.worker_loop
    @ (state: run_state, deadline_ns: i64) -> void
    + runs the scenario in a loop until the deadline, recording every sample
    # execution
    -> std.time.now_nanos
    -> std.time.sleep_nanos
    -> load_test.rate_gate
    -> load_test.record_sample
  load_test.run
    @ (state: run_state) -> metrics_snapshot
    + spawns worker_loop workers, waits for the deadline, and returns the final snapshot
    # execution
    -> std.sync.spawn
    -> std.time.now_nanos
    -> load_test.worker_loop
    -> load_test.snapshot
  load_test.format_tui_frame
    @ (snapshot: metrics_snapshot, width: i32, height: i32) -> string
    + renders a text frame with headline numbers and a small ASCII latency histogram
    # rendering
  load_test.live_frames
    @ (state: run_state, interval_ms: i32, emit: fn(string) -> void) -> void
    + calls emit with a fresh TUI frame at the given interval until the run ends
    # streaming
    -> std.time.now_nanos
    -> std.time.sleep_nanos
    -> load_test.snapshot
    -> load_test.format_tui_frame
