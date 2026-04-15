# Requirement: "a minimalist generic test assertion library"

A handful of assertion primitives that return a structured result rather than aborting. Callers decide what to do with failures.

std: (all units exist)

assert
  assert.equal_i64
    fn (want: i64, got: i64) -> result[void, string]
    + returns ok when the values are equal
    - returns a formatted mismatch message when they differ
    # assertion
  assert.equal_string
    fn (want: string, got: string) -> result[void, string]
    + returns ok when the strings are equal
    - returns a formatted mismatch showing both sides
    # assertion
  assert.is_true
    fn (cond: bool, label: string) -> result[void, string]
    + returns ok when cond is true
    - returns "<label>: expected true" when cond is false
    # assertion
  assert.contains
    fn (haystack: string, needle: string) -> result[void, string]
    + returns ok when needle is a substring of haystack
    - returns a message naming both sides when absent
    # assertion
