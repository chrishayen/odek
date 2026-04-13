# Requirement: "a benchmarking library with allocation profiling"

Runs a measured function repeatedly, records wall time and allocation counts, and summarizes the results.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time
  std.runtime
    std.runtime.alloc_count
      @ () -> i64
      + returns the total number of heap allocations since process start
      # runtime
    std.runtime.alloc_bytes
      @ () -> i64
      + returns the total bytes allocated since process start
      # runtime

bench
  bench.measure_once
    @ (label: string, body: fn() -> void) -> sample
    + runs body once and captures elapsed nanoseconds, alloc count delta, alloc bytes delta
    # measurement
    -> std.time.now_nanos
    -> std.runtime.alloc_count
    -> std.runtime.alloc_bytes
  bench.run
    @ (label: string, iterations: i32, body: fn() -> void) -> list[sample]
    + runs body `iterations` times and returns a sample per iteration
    ? iterations must be > 0; zero is a caller bug
    # measurement
    -> bench.measure_once
  bench.summarize
    @ (samples: list[sample]) -> bench_summary
    + returns min, max, mean, and median nanoseconds along with mean allocs
    # statistics
  bench.format_summary
    @ (label: string, summary: bench_summary) -> string
    + returns a human-readable one-line summary
    # formatting
