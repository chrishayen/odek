# Requirement: "an extendable assertion library with minimal boilerplate"

A small set of assertion runes that return structured results rather than panicking, so custom assertions compose easily.

std: (all units exist)

assertions
  assertions.equal
    @ (expected: string, actual: string) -> result[void, string]
    + returns ok when expected equals actual
    - returns a descriptive mismatch message otherwise
    # equality
  assertions.contains
    @ (haystack: string, needle: string) -> result[void, string]
    + returns ok when haystack contains needle
    - returns a mismatch message when needle is absent
    # substring
  assertions.all
    @ (checks: list[result[void, string]]) -> result[void, list[string]]
    + returns ok when every check passed
    - returns the list of failure messages otherwise
    # composition
  assertions.extend
    @ (name: string, predicate: fn(string) -> bool) -> assertion
    + returns a reusable assertion that runs the predicate with a labeled failure
    ? used to build domain-specific assertions without subclassing
    # extension
