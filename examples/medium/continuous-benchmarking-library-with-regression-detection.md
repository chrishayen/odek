# Requirement: "a continuous benchmarking library that detects performance regressions"

Ingests benchmark results over time and flags a measurement as a regression when it deviates significantly from the recent baseline.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

bencher
  bencher.new_store
    @ () -> store_state
    + creates an empty benchmark result store
    # construction
  bencher.record
    @ (store: store_state, benchmark_id: string, value: f64) -> store_state
    + appends a new measurement with the current timestamp
    # ingest
    -> std.time.now_millis
  bencher.baseline_stats
    @ (store: store_state, benchmark_id: string, window: i32) -> result[baseline_stats, string]
    + returns mean and stddev over the last `window` measurements
    - returns error when fewer than two measurements exist in the window
    # statistics
  bencher.detect_regression
    @ (store: store_state, benchmark_id: string, value: f64, sigma_threshold: f64) -> result[bool, string]
    + returns true when value exceeds baseline mean by more than sigma_threshold * stddev
    - returns error when a baseline cannot be computed
    # regression
  bencher.history
    @ (store: store_state, benchmark_id: string) -> list[tuple[i64, f64]]
    + returns all (timestamp, value) pairs for a benchmark in order
    # query
  bencher.clear
    @ (store: store_state, benchmark_id: string) -> store_state
    + removes all measurements for a benchmark
    # maintenance
