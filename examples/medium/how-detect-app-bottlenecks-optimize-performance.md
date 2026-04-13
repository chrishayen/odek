# Requirement: "an application performance profiler"

Collects timing samples from instrumented code blocks and produces hot-spot reports.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns a monotonic clock reading in nanoseconds
      # time

profiler
  profiler.new
    @ () -> profiler_state
    + creates an empty profiler
    # construction
  profiler.begin_span
    @ (p: profiler_state, name: string) -> tuple[i64, profiler_state]
    + starts a timed span and returns a span id
    # instrumentation
    -> std.time.now_nanos
  profiler.end_span
    @ (p: profiler_state, span_id: i64) -> profiler_state
    + records the elapsed time for a span against its name
    - ignores unknown span ids
    # instrumentation
    -> std.time.now_nanos
  profiler.record_sample
    @ (p: profiler_state, name: string, duration_nanos: i64) -> profiler_state
    + appends a single duration sample directly without using spans
    # instrumentation
  profiler.summary
    @ (p: profiler_state) -> list[tuple[string, i64, i64, i64, i64]]
    + returns (name, count, total_nanos, min_nanos, max_nanos) rows
    # reporting
  profiler.top_hotspots
    @ (p: profiler_state, n: i32) -> list[tuple[string, i64]]
    + returns the top n names by total time descending
    # reporting
  profiler.reset
    @ (p: profiler_state) -> profiler_state
    + discards all collected samples
    # lifecycle
