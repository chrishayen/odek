# Requirement: "an assertion library that extends a host test framework with expressive checks"

A small set of assertion primitives returning structured results the caller can report through its own test runner.

std: (all units exist)

assertions
  assertions.equal
    @ (expected: string, actual: string) -> result[void, string]
    + returns ok when both values are equal
    - returns error with a formatted diff message when they differ
    # equality
  assertions.not_equal
    @ (expected: string, actual: string) -> result[void, string]
    + returns ok when values differ
    - returns error when values are equal
    # equality
  assertions.contains
    @ (haystack: string, needle: string) -> result[void, string]
    + returns ok when haystack contains needle
    - returns error with haystack excerpt when needle is absent
    # substring
  assertions.true_check
    @ (value: bool, message: string) -> result[void, string]
    + returns ok when value is true
    - returns error carrying the provided message when value is false
    # boolean
  assertions.nil_check
    @ (value: optional[string]) -> result[void, string]
    + returns ok when the optional is empty
    - returns error when a value is present
    # optionality
  assertions.error_check
    @ (value: result[string, string]) -> result[void, string]
    + returns ok when the result carries an error
    - returns error when the result is a success
    # failure_expectation
