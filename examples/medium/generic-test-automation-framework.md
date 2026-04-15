# Requirement: "a generic test automation framework"

A library for registering named test cases, running them, and producing a structured report. No discovery, no DSL — callers register cases directly.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

test_framework
  test_framework.new_suite
    fn (name: string) -> suite_state
    + creates an empty suite with the given display name
    # construction
  test_framework.register
    fn (state: suite_state, case_name: string, body: string) -> suite_state
    + adds a named case whose body is an opaque handle the runner will invoke
    ? body is stored as a handle string the caller resolves via a runner callback
    # registration
  test_framework.run
    fn (state: suite_state) -> suite_report
    + executes every registered case in registration order and collects outcomes
    + records per-case duration in milliseconds
    # execution
    -> std.time.now_millis
  test_framework.record_pass
    fn (report: suite_report, case_name: string, duration_millis: i64) -> suite_report
    + marks a case as passed with its measured duration
    # reporting
  test_framework.record_fail
    fn (report: suite_report, case_name: string, message: string, duration_millis: i64) -> suite_report
    + marks a case as failed with a message and duration
    # reporting
  test_framework.summary
    fn (report: suite_report) -> map[string, i32]
    + returns counts keyed by "total", "passed", "failed"
    # reporting
  test_framework.failed_cases
    fn (report: suite_report) -> list[string]
    + returns the names of failed cases in registration order
    # reporting
