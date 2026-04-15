# Requirement: "a test framework with colorized terminal output"

A small assertion framework that collects results and renders a colorized summary. ANSI coloring is a thin std primitive.

std
  std.term
    std.term.colorize
      fn (text: string, color: string) -> string
      + wraps text in ANSI escape codes for the given color name
      + returns text unchanged when color is "none"
      # terminal
    std.term.supports_color
      fn () -> bool
      + returns true when stdout is a TTY and TERM is not "dumb"
      # terminal

test_framework
  test_framework.new_suite
    fn (name: string) -> suite_state
    + creates an empty test suite with the given name
    # construction
  test_framework.assert_equal
    fn (suite: suite_state, label: string, expected: string, actual: string) -> suite_state
    + records a passing result when expected equals actual
    - records a failing result with both values when they differ
    # assertion
  test_framework.assert_true
    fn (suite: suite_state, label: string, condition: bool) -> suite_state
    + records a passing result when condition is true
    - records a failing result when condition is false
    # assertion
  test_framework.run
    fn (suite: suite_state, body: func(suite_state) -> suite_state) -> suite_state
    + executes the body function and returns the final suite state
    ? body is the user's test function; framework does not discover tests
    # execution
  test_framework.render_report
    fn (suite: suite_state) -> string
    + returns a multi-line summary with green PASS and red FAIL markers
    + includes a count of passes and failures
    # reporting
    -> std.term.colorize
    -> std.term.supports_color
  test_framework.all_passed
    fn (suite: suite_state) -> bool
    + returns true when the suite has no recorded failures
    # reporting
