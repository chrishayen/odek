# Requirement: "pure generic functions over sequences"

A small toolkit of functional operations on lists. Each function is pure and returns a new list.

std: (all units exist)

slices
  slices.map
    fn (xs: list[T], f: fn(x: T) -> U) -> list[U]
    + returns a new list with the function applied to each element
    + returns an empty list for empty input
    # transform
  slices.filter
    fn (xs: list[T], pred: fn(x: T) -> bool) -> list[T]
    + returns a new list containing elements where the predicate holds
    + returns an empty list when no element matches
    # transform
  slices.reduce
    fn (xs: list[T], seed: U, f: fn(acc: U, x: T) -> U) -> U
    + returns the left-fold over the list using the seed
    + returns the seed unchanged for an empty list
    # aggregation
  slices.flatten
    fn (xs: list[list[T]]) -> list[T]
    + returns a single list concatenating all inner lists in order
    + returns an empty list for an empty input
    # transform
  slices.unique
    fn (xs: list[T]) -> list[T]
    + returns a new list with duplicates removed, preserving first-occurrence order
    ? element equality uses the host language's structural equality
    # transform
  slices.chunk
    fn (xs: list[T], size: i32) -> result[list[list[T]], string]
    + returns the list split into chunks of the given size
    + the final chunk may be shorter
    - returns error when size is zero or negative
    # transform
  slices.zip
    fn (xs: list[T], ys: list[U]) -> list[tuple[T, U]]
    + returns pairs until the shorter list is exhausted
    # transform
  slices.reverse
    fn (xs: list[T]) -> list[T]
    + returns the list with elements in reverse order
    # transform
