# Requirement: "a library of Haskell-inspired functions for lists"

A small set of pure list combinators. Only the ones with real content; host-language iteration is not a rune.

std: (all units exist)

hasgo
  hasgo.take
    fn (xs: list[i64], n: i32) -> list[i64]
    + returns the first n elements
    + returns the whole list when n is greater than the length
    + returns an empty list when n is zero or negative
    # slicing
  hasgo.drop
    fn (xs: list[i64], n: i32) -> list[i64]
    + returns the list without the first n elements
    + returns an empty list when n is greater than the length
    # slicing
  hasgo.take_while
    fn (xs: list[i64], pred: fn(i64) -> bool) -> list[i64]
    + returns the longest prefix whose elements satisfy pred
    - stops at the first element that fails pred
    # filtering
  hasgo.group_by
    fn (xs: list[i64], eq: fn(i64, i64) -> bool) -> list[list[i64]]
    + returns runs of adjacent elements that are equal under eq
    + returns an empty outer list for an empty input
    # grouping
  hasgo.scanl
    fn (xs: list[i64], seed: i64, step: fn(i64, i64) -> i64) -> list[i64]
    + returns each intermediate accumulator value, starting with seed
    + has length equal to xs length plus one
    # reduction
