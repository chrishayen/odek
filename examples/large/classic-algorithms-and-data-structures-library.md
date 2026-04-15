# Requirement: "a library of classic algorithms and data structures"

A curated set of data structures and algorithms that callers can compose. Each unit is self-contained.

std: (all units exist)

algorithms
  algorithms.quicksort
    fn (xs: list[i64]) -> list[i64]
    + returns xs sorted in ascending order
    + empty input returns an empty list
    # sorting
  algorithms.mergesort
    fn (xs: list[i64]) -> list[i64]
    + returns xs sorted in ascending order with a stable merge
    # sorting
  algorithms.binary_search
    fn (xs: list[i64], target: i64) -> i32
    + returns the index of target in a sorted list
    - returns -1 when the target is not present
    # searching
  algorithms.heap_push
    fn (heap: list[i64], value: i64) -> list[i64]
    + returns a new min-heap with value inserted
    # heap
  algorithms.heap_pop
    fn (heap: list[i64]) -> tuple[optional[i64], list[i64]]
    + returns (min_value, new_heap)
    - returns (none, heap) when the heap is empty
    # heap
  algorithms.bst_insert
    fn (tree: bst_state, key: i64, value: string) -> bst_state
    + returns a tree containing the key-value pair
    + updates the value when the key already exists
    # binary_search_tree
  algorithms.bst_lookup
    fn (tree: bst_state, key: i64) -> optional[string]
    + returns the associated value when the key is present
    # binary_search_tree
  algorithms.bst_delete
    fn (tree: bst_state, key: i64) -> bst_state
    + returns a tree with the key removed
    # binary_search_tree
  algorithms.bfs
    fn (graph: graph_state, source: i32) -> list[i32]
    + returns nodes in breadth-first visit order
    # graph_search
  algorithms.dfs
    fn (graph: graph_state, source: i32) -> list[i32]
    + returns nodes in depth-first visit order
    # graph_search
  algorithms.dijkstra
    fn (graph: graph_state, source: i32) -> map[i32, f64]
    + returns shortest distances from source to every reachable node
    - distances are +infinity for unreachable nodes
    ? edge weights must be non-negative
    # shortest_paths
  algorithms.kruskal_mst
    fn (graph: graph_state) -> list[edge]
    + returns edges forming a minimum spanning tree
    - returns empty list when the graph is empty
    # spanning_trees
  algorithms.union_find_new
    fn (n: i32) -> uf_state
    + returns a disjoint-set structure with n singleton sets
    # union_find
  algorithms.union_find_union
    fn (uf: uf_state, a: i32, b: i32) -> uf_state
    + merges the sets containing a and b using union by rank
    # union_find
  algorithms.union_find_find
    fn (uf: uf_state, a: i32) -> tuple[i32, uf_state]
    + returns (root, compressed_state)
    # union_find
  algorithms.longest_common_subseq
    fn (a: string, b: string) -> string
    + returns the longest common subsequence of a and b
    + returns empty string when either input is empty
    # dynamic_programming
  algorithms.knapsack_01
    fn (weights: list[i32], values: list[i32], capacity: i32) -> i32
    + returns the maximum total value fitting under the capacity
    # dynamic_programming
