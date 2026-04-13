# Requirement: "a library of convenient iterator combinators"

Combinators that transform and consume lazy iterators. An iterator is an opaque handle that yields values one at a time.

std: (all units exist)

iter_combinators
  iter_combinators.from_list
    @ (values: list[i64]) -> iterator
    + returns an iterator over the list's elements in order
    # construction
  iter_combinators.map
    @ (source: iterator, fn: function[i64, i64]) -> iterator
    + returns an iterator that applies fn to each element lazily
    # transform
  iter_combinators.filter
    @ (source: iterator, predicate: function[i64, bool]) -> iterator
    + returns an iterator that yields only elements satisfying predicate
    # transform
  iter_combinators.take
    @ (source: iterator, n: i32) -> iterator
    + returns an iterator that yields at most the first n elements
    # transform
  iter_combinators.fold
    @ (source: iterator, initial: i64, fn: function[i64, i64, i64]) -> i64
    + accumulates over the iterator starting from initial
    + returns initial for an empty iterator
    # consumption
