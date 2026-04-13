# Requirement: "a concurrent sorting library"

Partitions the input and sorts chunks in parallel, merging the results. The caller provides a comparator.

std
  std.concurrency
    std.concurrency.parallel_map
      @ (chunks: list[list[i64]], f: fn(list[i64]) -> list[i64]) -> list[list[i64]]
      + applies f to each chunk in parallel and returns the results in order
      # parallelism

parallel_sort
  parallel_sort.sort
    @ (items: list[i64], workers: i32) -> list[i64]
    + returns items sorted ascending using up to workers threads
    + returns items unchanged when its length is below the parallel threshold
    ? falls back to a sequential sort when workers <= 1
    # sort
    -> std.concurrency.parallel_map
  parallel_sort.sort_by
    @ (items: list[i64], workers: i32, less: fn(i64, i64) -> bool) -> list[i64]
    + sorts using the caller's comparator
    # sort
    -> std.concurrency.parallel_map
  parallel_sort.merge_sorted
    @ (a: list[i64], b: list[i64], less: fn(i64, i64) -> bool) -> list[i64]
    + merges two pre-sorted lists into one
    # merge
