# Requirement: "a micro-benchmarking library for measuring the speed of small code constructs"

Runs a caller-supplied closure many times and reports timing statistics. A thin std clock primitive keeps benchmarks deterministic under test.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time

benchmark
  benchmark.run
    fn (iterations: i32, body: fn() -> void) -> benchmark_result
    + invokes body the given number of times and records per-iteration nanoseconds
    - returns a zero-iteration result when iterations <= 0
    # measurement
    -> std.time.now_nanos
  benchmark.summarize
    fn (result: benchmark_result) -> benchmark_summary
    + returns min, max, mean, and total elapsed nanoseconds
    + reports iterations_per_second as f64
    # statistics
