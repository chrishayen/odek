# Requirement: "a reactive observable stream library"

An observable is a value you can subscribe to, transform, and combine. The library is pull-based so tests don't need a scheduler.

std: (all units exist)

rx
  rx.from_list
    @ (values: list[T]) -> observable[T]
    + returns an observable that emits each value in order and completes
    # source
  rx.from_fn
    @ (producer: func() -> optional[T]) -> observable[T]
    + returns an observable driven by the producer; completes when producer returns none
    # source
  rx.map
    @ (source: observable[T], f: func(T) -> U) -> observable[U]
    + returns an observable that applies f to each emission
    # transform
  rx.filter
    @ (source: observable[T], pred: func(T) -> bool) -> observable[T]
    + returns an observable that drops emissions for which pred is false
    # transform
  rx.merge
    @ (left: observable[T], right: observable[T]) -> observable[T]
    + returns an observable that interleaves emissions from both sources
    + completes after both sources complete
    # combine
  rx.scan
    @ (source: observable[T], initial: U, f: func(U, T) -> U) -> observable[U]
    + emits running fold values, starting with initial
    # transform
  rx.take
    @ (source: observable[T], n: i32) -> observable[T]
    + emits at most n values from source then completes
    - completes immediately when n is zero or negative
    # transform
  rx.collect
    @ (source: observable[T]) -> list[T]
    + drains source and returns all emitted values in order
    # sink
  rx.subscribe
    @ (source: observable[T], on_next: func(T) -> void, on_complete: func() -> void) -> void
    + drives source to completion, invoking callbacks for each emission
    # sink
