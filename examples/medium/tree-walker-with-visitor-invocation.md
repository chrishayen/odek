# Requirement: "a library that walks a structured value tree and invokes a visitor at each node"

A generic recursive walker for structured values (scalars, lists, maps). The caller provides a visitor function.

std: (all units exist)

walker
  walker.walk
    fn (root: dynamic_value, visit: visit_fn) -> void
    + invokes visit at the root, then recurses into children
    + visits list elements in index order
    + visits map entries in insertion order
    # traversal
  walker.walk_with_path
    fn (root: dynamic_value, visit: visit_with_path_fn) -> void
    + passes a path breadcrumb (e.g. ".items[2].name") to each visit call
    # traversal
  walker.find_all
    fn (root: dynamic_value, predicate: predicate_fn) -> list[path_and_value]
    + returns every node whose predicate returns true, with its path
    # query
  walker.transform
    fn (root: dynamic_value, transform: transform_fn) -> dynamic_value
    + returns a new tree with each node replaced by the transform's output
    + does not mutate the input
    # rewriting
  walker.depth
    fn (root: dynamic_value) -> i32
    + returns the maximum nesting depth of the tree
    + returns 1 for a scalar
    # inspection
