# Requirement: "a database benchmarking tool with support for several databases and scripts"

A benchmark runner that executes scripted workloads against a pluggable database driver and reports latency/throughput statistics.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time
  std.math
    std.math.percentile
      @ (samples: list[i64], p: f64) -> i64
      + returns the p-th percentile of the samples
      - returns 0 when samples is empty
      # statistics

bench
  bench.load_script
    @ (source: string) -> result[list[string], string]
    + parses a script into a list of statements separated by semicolons
    - returns error on unterminated quoted strings
    # scripting
  bench.new_run
    @ (driver_name: string, script: list[string], concurrency: i32, iterations: i32) -> bench_run
    + builds a run configuration bound to a registered driver
    - returns a run whose execute will fail if the driver is not registered
    # construction
  bench.register_driver
    @ (name: string, executor: driver_fn) -> void
    + associates a name with a statement-executing function
    # driver_registry
  bench.execute
    @ (run: bench_run) -> result[bench_report, string]
    + runs the script across workers and collects per-statement latencies
    - returns error when concurrency is less than one
    # execution
    -> std.time.now_nanos
  bench.summarize
    @ (report: bench_report) -> bench_summary
    + computes throughput and p50/p95/p99 latency for the run
    # reporting
    -> std.math.percentile
