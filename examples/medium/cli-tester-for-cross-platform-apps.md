# Requirement: "a library for testing command-line applications across platforms"

Declarative test cases describe a command to run, expected exit code, and expected stdout/stderr matchers. The runner executes them and reports pass/fail.

std
  std.process
    std.process.run
      fn (command: string, args: list[string]) -> result[process_result, string]
      + runs the command and returns exit code, stdout, and stderr
      - returns error when the command cannot be launched
      # process
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem

cli_test
  cli_test.load_suite
    fn (path: string) -> result[list[test_case], string]
    + returns the parsed list of test cases from a suite file
    - returns error when the file cannot be read or has invalid structure
    # loading
    -> std.fs.read_all
  cli_test.match_output
    fn (actual: string, matcher: output_matcher) -> bool
    + returns true when actual satisfies an exact, contains, or regex matcher
    # matching
  cli_test.run_case
    fn (case: test_case) -> test_result
    + runs the case's command and returns a result carrying pass/fail and any mismatched expectations
    + fails when exit code, stdout, or stderr do not match the expectations
    # execution
    -> std.process.run
  cli_test.run_suite
    fn (cases: list[test_case]) -> suite_report
    + runs every case and returns a report with totals and per-case results
    # execution
