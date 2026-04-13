# Requirement: "a collection of functional helpers for sequences and functions"

Classic functional combinators: map, filter, reduce, partition, group_by, and function composition.

std: (all units exist)

funcy
  funcy.map
    @ (items: list[i64], f: fn(i64) -> i64) -> list[i64]
    + returns a list with f applied to each element
    + returns an empty list when input is empty
    # transform
  funcy.filter
    @ (items: list[i64], pred: fn(i64) -> bool) -> list[i64]
    + returns elements where pred returns true, preserving order
    # selection
  funcy.reduce
    @ (items: list[i64], init: i64, f: fn(i64, i64) -> i64) -> i64
    + folds items left-to-right starting from init
    + returns init when items is empty
    # folding
  funcy.partition
    @ (items: list[i64], pred: fn(i64) -> bool) -> tuple[list[i64], list[i64]]
    + returns (matches, non_matches) with original order preserved in each
    # partitioning
  funcy.group_by
    @ (items: list[i64], key_fn: fn(i64) -> string) -> map[string, list[i64]]
    + groups items by their computed key
    # grouping
  funcy.compose
    @ (f: fn(i64) -> i64, g: fn(i64) -> i64) -> fn(i64) -> i64
    + returns a function equivalent to applying g then f
    ? right-to-left composition like mathematical notation
    # composition
