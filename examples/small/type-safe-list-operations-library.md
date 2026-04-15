# Requirement: "a library of type-safe list operations"

Generic list helpers beyond what the host language provides: chunking, unique, partition, zip.

std: (all units exist)

list_ops
  list_ops.chunk
    fn (xs: list[i32], size: i32) -> list[list[i32]]
    + splits xs into sublists of at most size
    + returns an empty list when xs is empty
    - returns an empty list when size is not positive
    # chunking
  list_ops.unique
    fn (xs: list[i32]) -> list[i32]
    + returns xs with duplicates removed, preserving first-seen order
    # deduplication
  list_ops.partition
    fn (xs: list[i32], predicate: i32 -> bool) -> tuple[list[i32], list[i32]]
    + returns (matching, non_matching) preserving order within each
    # partitioning
  list_ops.zip
    fn (xs: list[i32], ys: list[string]) -> list[tuple[i32, string]]
    + pairs elements by position, truncating to the shorter input
    # zipping
  list_ops.group_by
    fn (xs: list[i32], key: i32 -> string) -> map[string, list[i32]]
    + groups elements by the result of key, preserving insertion order per group
    # grouping
