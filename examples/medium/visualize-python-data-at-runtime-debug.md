# Requirement: "a runtime data visualizer for debugging references, mutability, and aliasing"

Snapshots a value graph at runtime and produces a diagram description showing objects, their fields, and shared references (aliasing).

std
  std.collections
    std.collections.hashmap_new
      @ () -> map[i64, object_snapshot]
      + returns an empty object id map
      # collections
    std.collections.hashmap_put
      @ (m: map[i64, object_snapshot], key: i64, value: object_snapshot) -> map[i64, object_snapshot]
      + returns a new map with key bound to value
      # collections
    std.collections.hashmap_contains
      @ (m: map[i64, object_snapshot], key: i64) -> bool
      + returns true when key is present
      # collections

memviz
  memviz.snapshot
    @ (root: dynamic_value) -> object_graph
    + traverses root recording each distinct object by its identity
    + revisits of already-seen objects are recorded as alias edges, not copies
    # capture
    -> std.collections.hashmap_new
    -> std.collections.hashmap_put
    -> std.collections.hashmap_contains
  memviz.object_id
    @ (val: dynamic_value) -> i64
    + returns a stable runtime identity token for val
    ? scalar values produce synthetic ids to keep the graph total
    # identity
  memviz.find_aliases
    @ (graph: object_graph) -> list[alias_edge]
    + returns every pair of fields that point at the same object id
    + returns [] when no sharing exists
    # analysis
  memviz.diff
    @ (before: object_graph, after: object_graph) -> list[mutation]
    + returns the set of field changes between two snapshots
    + marks added, removed, and modified fields separately
    # analysis
  memviz.to_diagram
    @ (graph: object_graph) -> string
    + produces a diagram description with nodes for objects and edges for references
    + marks alias edges distinctly from owning edges
    # rendering
