# Requirement: "an assertion library with detailed structural diffs and human-readable failure messages"

The project exposes a small set of assertions; diff rendering lives in its own rune and leans on a generic value model.

std: (all units exist)

assert
  assert.equal
    fn (expected: value, actual: value) -> result[void, string]
    + returns ok when expected and actual are structurally equal
    - returns error with a rendered diff when they differ
    # equality
  assert.not_equal
    fn (expected: value, actual: value) -> result[void, string]
    + returns ok when the values differ
    - returns error when the values are structurally equal
    # inequality
  assert.contains
    fn (haystack: value, needle: value) -> result[void, string]
    + returns ok when a list or map contains the needle
    - returns error with the missing element in the message
    # containment
  assert.is_nil
    fn (v: value) -> result[void, string]
    + returns ok for a nil/empty optional value
    - returns error describing the non-nil value
    # nullability
  assert.diff_values
    fn (expected: value, actual: value) -> string
    + returns a multi-line diff with paths like ".user.name" marking changed fields
    + aligned +/- lines for each differing leaf
    + returns "" when the values are equal
    # diff_rendering
  assert.render_value
    fn (v: value) -> string
    + returns a stable, human-readable rendering of any structural value
    + nested maps and lists are indented consistently
    # rendering
  assert.format_failure
    fn (message: string, diff: string) -> string
    + returns a multi-line failure message with the header and the diff
    # formatting
