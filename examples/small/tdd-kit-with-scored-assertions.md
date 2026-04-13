# Requirement: "a library supporting test-driven development exercises with scored assertions"

A micro test harness: register a case, run it against a candidate function, report pass/fail.

std: (all units exist)

tddkit
  tddkit.new_suite
    @ (name: string) -> test_suite
    + creates an empty suite with a name
    # construction
  tddkit.add_case
    @ (suite: test_suite, description: string, input: string, expected: string) -> test_suite
    + appends a case with its expected output
    # cases
  tddkit.run
    @ (suite: test_suite, candidate: func_handle) -> test_report
    + runs each case through the candidate and records pass/fail
    # execution
  tddkit.summary
    @ (report: test_report) -> tuple[i32, i32]
    + returns (passed, failed) counts
    # reporting
