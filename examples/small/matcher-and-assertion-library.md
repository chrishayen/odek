# Requirement: "a matcher and assertion library"

Fluent matchers that each return an optional error message. No runner — callers plug these into any test framework.

std: (all units exist)

matchers
  matchers.expect_equal
    fn (actual: string, expected: string) -> optional[string]
    + returns none when values are equal
    - returns "expected <expected>, got <actual>" when values differ
    # equality
  matchers.expect_contains
    fn (haystack: string, needle: string) -> optional[string]
    + returns none when haystack contains needle as a substring
    - returns a descriptive error when needle is not found
    # substring
  matchers.expect_length
    fn (items: list[string], expected: i32) -> optional[string]
    + returns none when the list length matches
    - returns an error describing the observed length when it does not match
    # length
  matchers.expect_true
    fn (value: bool, label: string) -> optional[string]
    + returns none when value is true
    - returns "expected <label> to be true" when value is false
    # truthiness
