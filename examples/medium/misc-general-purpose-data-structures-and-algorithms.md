# Requirement: "a small collection of general-purpose data structures and algorithms"

A handful of widely-useful primitives: a min-heap, a union-find, and one graph traversal.

std: (all units exist)

misc_ds
  misc_ds.heap_new
    fn () -> heap_state
    + returns an empty min-heap of integers
    # construction
  misc_ds.heap_push
    fn (state: heap_state, value: i32) -> heap_state
    + returns new state with value inserted in heap order
    # mutation
  misc_ds.heap_pop_min
    fn (state: heap_state) -> tuple[optional[i32], heap_state]
    + returns (smallest, new_state) when non-empty
    - returns (empty, unchanged_state) when empty
    # mutation
  misc_ds.union_find_new
    fn (n: i32) -> union_find_state
    + returns a forest with n disjoint singletons numbered 0..n-1
    # construction
  misc_ds.union_find_union
    fn (state: union_find_state, a: i32, b: i32) -> union_find_state
    + merges the sets containing a and b using rank
    # mutation
  misc_ds.union_find_find
    fn (state: union_find_state, x: i32) -> tuple[i32, union_find_state]
    + returns (root, new_state) and applies path compression
    # lookup
  misc_ds.bfs
    fn (edges: map[i32, list[i32]], start: i32) -> list[i32]
    + returns vertices in breadth-first visit order from start
    + returns just [start] when the vertex has no edges
    # graph_traversal
