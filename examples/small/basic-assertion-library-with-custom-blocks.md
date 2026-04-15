# Requirement: "a basic assertion library with building blocks for custom assertions"

Provides a handful of equality and relational checks that return a structured outcome rather than throwing. Callers compose these into their own higher-level assertions.

std: (all units exist)

assert
  assert.equal
    fn (expected: string, actual: string) -> assertion_result
    + returns ok when expected equals actual
    - returns a failing result with a diff message otherwise
    # equality
  assert.not_equal
    fn (expected: string, actual: string) -> assertion_result
    + returns ok when expected differs from actual
    - returns a failing result when they are equal
    # equality
  assert.contains
    fn (haystack: string, needle: string) -> assertion_result
    + returns ok when haystack contains needle
    - returns a failing result with both values otherwise
    # substring
  assert.list_equal
    fn (expected: list[string], actual: list[string]) -> assertion_result
    + returns ok when both lists have the same length and elements in order
    - returns a failing result describing the first divergence
    # collection
  assert.compose
    fn (name: string, checks: list[assertion_result]) -> assertion_result
    + returns ok when every check passes
    - returns a failing result that aggregates the failing messages under name
    # composition
