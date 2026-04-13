# Requirement: "a benchmark harness for database query drivers"

A library that times named query runs against pluggable driver adapters and produces aggregate results.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns the current monotonic time in nanoseconds
      # time

db_bench
  db_bench.new_suite
    @ (name: string) -> suite_state
    + creates an empty benchmark suite
    # construction
  db_bench.add_case
    @ (state: suite_state, name: string, driver_name: string, query: string, iterations: i32) -> suite_state
    + registers a benchmark case
    # registration
  db_bench.record_run
    @ (state: suite_state, case_name: string, elapsed_nanos: i64) -> suite_state
    + appends a timing sample for a case
    # recording
  db_bench.time_block
    @ (state: suite_state, case_name: string, start_nanos: i64) -> suite_state
    + records elapsed time from a start marker to now
    # recording
    -> std.time.now_nanos
  db_bench.summarize
    @ (state: suite_state) -> list[case_summary]
    + returns per-case min, max, mean, and p95 timings in nanoseconds
    ? requires at least one sample per case; cases with no samples are omitted
    # reporting
