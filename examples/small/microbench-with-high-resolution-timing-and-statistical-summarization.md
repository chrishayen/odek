# Requirement: "a micro-benchmark library with high-resolution timing and statistical summaries"

Runs a function repeatedly with a warmup, collects per-iteration nanosecond timings, then reports statistics.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns the current monotonic time in nanoseconds
      # time

microbench
  microbench.run
    @ (name: string, body: fn() -> void, iterations: i32, warmup: i32) -> bench_result
    + runs warmup iterations, then records timings for iterations and returns the collected samples
    ? iterations must be at least 1; warmup may be 0
    # measurement
    -> std.time.now_nanos
  microbench.statistics
    @ (result: bench_result) -> bench_stats
    + returns mean, standard deviation, min, max, and p50/p95/p99 in nanoseconds
    # statistics
  microbench.is_significant
    @ (a: bench_stats, b: bench_stats) -> bool
    + returns true when the 95% confidence intervals of a and b do not overlap
    # comparison
  microbench.format
    @ (name: string, stats: bench_stats) -> string
    + returns a human-readable one-line summary like "name: 123 ns/op +/- 5 ns"
    # reporting
