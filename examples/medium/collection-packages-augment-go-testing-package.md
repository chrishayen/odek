# Requirement: "a collection of packages to augment a unit testing framework with common assertion patterns"

Assertion helpers that report rich failure messages and structure comparisons.

std
  std.fmt
    std.fmt.sprintf
      @ (template: string, args: list[string]) -> string
      + substitutes positional placeholders
      # formatting
  std.reflect
    std.reflect.deep_equal
      @ (a: bytes, b: bytes) -> bool
      + returns true when two opaque values have identical structure
      # reflection

testutil
  testutil.assert_equal
    @ (actual: string, expected: string) -> result[void, string]
    + returns ok when values match
    - returns descriptive error when they differ
    # assertion
    -> std.fmt.sprintf
  testutil.assert_nil
    @ (value: optional[string]) -> result[void, string]
    + returns ok when value is absent
    - returns error with a formatted diff when value is present
    # assertion
  testutil.assert_contains
    @ (haystack: string, needle: string) -> result[void, string]
    + returns ok when haystack contains needle
    - returns error listing haystack and needle when missing
    # assertion
  testutil.assert_deep_equal
    @ (actual: bytes, expected: bytes) -> result[void, string]
    + returns ok when structures match
    - returns error with a path to the first difference
    # assertion
    -> std.reflect.deep_equal
  testutil.assert_panics
    @ (thunk: bytes) -> result[void, string]
    + returns ok when the thunk aborts
    - returns error when the thunk completes normally
    ? thunk is an opaque deferred call captured by the caller
    # assertion
