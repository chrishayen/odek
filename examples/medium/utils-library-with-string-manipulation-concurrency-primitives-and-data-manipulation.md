# Requirement: "a grab-bag utility library covering string manipulation, concurrency primitives, and data manipulation"

A small curated surface. Real utilities only; no filler.

std: (all units exist)

utils
  utils.string_reverse
    @ (s: string) -> string
    + returns the string with codepoints in reverse order
    + handles empty string
    # strings
  utils.string_contains_any
    @ (s: string, needles: list[string]) -> bool
    + returns true when any needle appears in s
    - returns false when needles is empty
    # strings
  utils.unique
    @ (items: list[string]) -> list[string]
    + returns items with duplicates removed, preserving first occurrence
    # collections
  utils.chunk
    @ (items: list[string], size: i32) -> list[list[string]]
    + splits a list into fixed-size chunks
    - returns empty when size is non-positive
    # collections
  utils.parallel_map
    @ (items: list[string], workers: i32, fn_id: i32) -> list[string]
    + applies a registered transform function across items using the given worker count
    ? fn_id refers to a transform registered out-of-band by the caller
    # concurrency
  utils.debounce
    @ (interval_ms: i64, fn_id: i32) -> i32
    + returns a handle that forwards only the last call within each interval
    # concurrency
