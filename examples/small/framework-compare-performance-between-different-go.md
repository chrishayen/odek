# Requirement: "a framework to compare the performance of a workload across different runtime versions"

The framework records timing samples tagged by version, then produces a comparison report showing mean, min, max, and relative speedup between versions.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time

benchcompare
  benchcompare.new
    @ () -> suite_state
    + creates an empty suite with no recorded samples
    # construction
  benchcompare.record
    @ (state: suite_state, version_label: string, benchmark_name: string, duration_nanos: i64) -> suite_state
    + appends a sample under (version_label, benchmark_name)
    - returns unchanged state when duration_nanos is negative
    # data_capture
  benchcompare.time_call
    @ (state: suite_state, version_label: string, benchmark_name: string, elapsed_fn_result: i64) -> suite_state
    + records a sample whose duration was measured by the caller against a monotonic clock
    # data_capture
    -> std.time.now_nanos
  benchcompare.summary
    @ (state: suite_state, version_label: string, benchmark_name: string) -> optional[stat_summary]
    + returns count, mean, min, and max of recorded samples
    - returns none when no samples exist for that key
    # statistics
  benchcompare.compare
    @ (state: suite_state, benchmark_name: string, baseline: string, candidate: string) -> result[f64, string]
    + returns candidate_mean / baseline_mean as the relative slowdown factor
    - returns error when either version has zero samples for benchmark_name
    # comparison
