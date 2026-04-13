# Requirement: "a benchmarking library that times repeated executions of a command"

Runs a callable many times and reports summary statistics.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic timestamp in nanoseconds
      # time
  std.process
    std.process.run
      @ (argv: list[string]) -> result[i32, string]
      + runs the command and returns its exit status
      - returns error when the program cannot be spawned
      # process

benchmark
  benchmark.time_once
    @ (argv: list[string]) -> result[i64, string]
    + returns the wall-clock duration in nanoseconds of running argv once
    - returns error when the process fails to start
    # timing
    -> std.time.now_nanos
    -> std.process.run
  benchmark.run_samples
    @ (argv: list[string], samples: i32) -> result[list[i64], string]
    + returns one duration per sample
    - returns error when samples is less than 1
    # sampling
    -> benchmark.time_once
  benchmark.summarize
    @ (durations: list[i64]) -> benchmark_summary
    + returns min, max, mean, and standard deviation of the samples
    ? standard deviation uses the sample (n-1) formula
    # statistics
