# Requirement: "a static assertion library for checking invariants"

Since the host may or may not support compile-time evaluation, the library exposes assertions as pure predicates on a spec value; the caller wires them into their build or test step. Failures return a reason rather than panicking.

std: (all units exist)

static_assert
  static_assert.size_eq
    @ (type_size_bytes: i32, expected: i32) -> result[void, string]
    + succeeds when type_size_bytes equals expected
    - returns an error message "expected size N, got M" otherwise
    # size
  static_assert.size_le
    @ (type_size_bytes: i32, max: i32) -> result[void, string]
    + succeeds when type_size_bytes is at most max
    - returns error otherwise
    # size
  static_assert.const_eq
    @ (actual: i64, expected: i64, label: string) -> result[void, string]
    + succeeds when actual equals expected
    - returns an error message including label otherwise
    # equality
  static_assert.all
    @ (checks: list[result[void, string]]) -> result[void, string]
    + succeeds when every check succeeds
    - returns the first failing check's error, prefixed with its index
    # composition
