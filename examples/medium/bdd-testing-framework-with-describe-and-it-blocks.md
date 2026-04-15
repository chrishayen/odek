# Requirement: "a BDD-style testing framework with describe and it blocks"

Users build a tree of suites and specs, then run it to get a structured report.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

tests
  tests.new_suite
    fn (name: string) -> suite
    + creates a top-level suite with the given name
    # construction
  tests.describe
    fn (parent: suite, name: string, body: fn(suite) -> void) -> suite
    + adds a nested suite and runs body with it so body can register specs
    # grouping
  tests.it
    fn (parent: suite, name: string, body: fn() -> result[void, string]) -> suite
    + registers a spec with a pass/fail body
    # specification
  tests.before_each
    fn (parent: suite, hook: fn() -> void) -> suite
    + registers a hook run before every spec in this suite and its descendants
    # lifecycle
  tests.after_each
    fn (parent: suite, hook: fn() -> void) -> suite
    + registers a hook run after every spec in this suite and its descendants
    # lifecycle
  tests.run
    fn (root: suite) -> run_report
    + executes all specs depth-first, invoking before_each and after_each hooks around each spec
    + captures per-spec pass/fail status, error message, and elapsed time
    # execution
    -> std.time.now_millis
  tests.format_report
    fn (report: run_report) -> string
    + renders a hierarchical text summary with counts of passed, failed, and skipped specs
    # reporting
