# Requirement: "an iterator with map and reduce"

A lazy iterator abstraction with map and reduce combinators. Iteration state is opaque so backing sources can vary.

std: (all units exist)

iter
  iter.from_list
    fn (items: list[i64]) -> iter_state
    + wraps a list as an iterator that yields each element in order
    # construction
  iter.next
    fn (state: iter_state) -> tuple[optional[i64], iter_state]
    + returns the next element and the advanced state
    - returns none when the iterator is exhausted
    # traversal
  iter.map
    fn (state: iter_state, f: fn(i64) -> i64) -> iter_state
    + returns a new iterator that applies f to each yielded element lazily
    # transformation
  iter.reduce
    fn (state: iter_state, initial: i64, f: fn(i64, i64) -> i64) -> i64
    + folds the iterator left-to-right starting from initial
    + returns initial when the iterator is empty
    # aggregation
