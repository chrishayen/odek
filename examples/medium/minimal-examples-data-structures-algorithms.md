# Requirement: "a library of minimal data structures and algorithms"

A tight collection of general-purpose building blocks: one container, one sort, one search.

std: (all units exist)

algos
  algos.stack_new
    @ () -> stack_state
    + returns an empty stack
    # construction
  algos.stack_push
    @ (state: stack_state, value: i32) -> stack_state
    + returns a new state with value on top
    # mutation
  algos.stack_pop
    @ (state: stack_state) -> tuple[optional[i32], stack_state]
    + returns (top_value, new_state) when non-empty
    - returns (empty, unchanged_state) when empty
    # mutation
  algos.quicksort
    @ (values: list[i32]) -> list[i32]
    + returns an ascending-sorted copy
    + returns an empty list when input is empty
    ? uses median-of-three pivot selection
    # sorting
  algos.binary_search
    @ (sorted: list[i32], target: i32) -> optional[i32]
    + returns the index when target is present
    - returns empty when target is absent
    ? assumes the input is already sorted ascending
    # search
