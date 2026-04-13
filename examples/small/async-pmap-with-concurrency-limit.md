# Requirement: "a library for mapping an async function over items with a concurrency limit"

Runs an async mapper across a list with a bounded number of in-flight tasks.

std: (all units exist)

pmap
  pmap.map
    @ (items: list[T], mapper: fn(T) -> result[U, string], concurrency: i32) -> result[list[U], string]
    + returns results in the original input order regardless of completion order
    + never has more than concurrency mapper calls running at once
    - returns the first mapper error and cancels in-flight work
    ? concurrency of zero or less is treated as one
    # concurrency
  pmap.map_settled
    @ (items: list[T], mapper: fn(T) -> result[U, string], concurrency: i32) -> list[result[U, string]]
    + runs every mapper call to completion even when some fail
    + result order matches input order
    # concurrency
