# Requirement: "a benchmark harness for common basic operations"

A small harness that times a user-provided operation repeatedly and reports per-op cost.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns the current monotonic time in nanoseconds
      # time

bench
  bench.run
    fn (name: string, iterations: i32, op: fn() -> void) -> bench_result
    + runs the op the given number of iterations and records total elapsed nanoseconds
    + computes nanoseconds per iteration
    - returns a zero-iteration result when iterations <= 0
    # measurement
    -> std.time.now_nanos
  bench.format_result
    fn (result: bench_result) -> string
    + returns a line like "name: N ops, X ns/op"
    # reporting
