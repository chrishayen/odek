# Requirement: "a slice utilities library"

Common operations on homogeneous lists that the host language does not provide as primitives.

std: (all units exist)

slices
  slices.chunk
    fn (xs: list[i64], size: i32) -> list[list[i64]]
    + splits xs into consecutive chunks of at most size elements
    - returns empty list when size <= 0
    # slicing
  slices.unique
    fn (xs: list[i64]) -> list[i64]
    + returns the first occurrence of each value, preserving order
    # deduplication
  slices.index_of
    fn (xs: list[i64], needle: i64) -> optional[i32]
    + returns the zero-based index of the first match
    - returns none when not present
    # search
  slices.partition
    fn (xs: list[i64], predicate: fn(i64) -> bool) -> tuple[list[i64], list[i64]]
    + returns (matching, non_matching) preserving original order
    # partitioning
  slices.rotate_left
    fn (xs: list[i64], n: i32) -> list[i64]
    + rotates elements left by n positions, wrapping around
    + n is taken modulo len(xs)
    - returns xs unchanged when xs is empty
    # rotation
