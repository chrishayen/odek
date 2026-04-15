# Requirement: "a load testing library that drives scripted scenarios and aggregates metrics"

Runs user-supplied request scenarios at a target rate and records latencies, successes, and failures.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns monotonic nanoseconds since an arbitrary epoch
      # time

load
  load.new_scenario
    fn (name: string, request: fn() -> result[void, string]) -> scenario
    + creates a scenario with a name and a callable that performs one iteration
    # construction
  load.new_plan
    fn () -> plan_state
    + creates an empty plan with no scenarios
    # construction
  load.stage
    fn (p: plan_state, s: scenario, rps: f64, duration_seconds: f64) -> plan_state
    + appends a stage that drives s at the given rate for the given duration
    # planning
  load.run
    fn (p: plan_state) -> run_report
    + executes each stage sequentially, recording the outcome and latency of every iteration
    # execution
    -> std.time.now_nanos
  load.summary
    fn (report: run_report) -> metric_summary
    + computes count, success count, failure count, and latency p50/p90/p99 per scenario
    # metrics
  load.format_summary
    fn (summary: metric_summary) -> string
    + renders a human-readable summary table
    # reporting
