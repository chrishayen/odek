# Requirement: "colorful test assertion helpers"

Assertions return a structured result; a separate renderer produces ANSI-colored output.

std
  std.term
    std.term.ansi_color
      fn (text: string, color: string) -> string
      + wraps text in the ANSI escape codes for the named color
      # terminal

colorful_test
  colorful_test.assert_equal
    fn (expected: string, actual: string) -> assertion_result
    + returns a passing result when expected equals actual
    - returns a failing result with both values recorded when they differ
    # assertion
  colorful_test.assert_true
    fn (label: string, value: bool) -> assertion_result
    + returns a passing result when value is true
    - returns a failing result labeled with the given name when false
    # assertion
  colorful_test.render
    fn (result: assertion_result) -> string
    + returns a green pass line or a red fail line with expected and actual values
    # rendering
    -> std.term.ansi_color
