# Requirement: "a library for topological sorting and pruning of dag graphs"

A directed acyclic graph utility with two core operations: topological sort and subgraph pruning.

std: (all units exist)

dag
  dag.new
    fn () -> dag_state
    + returns an empty graph
    # construction
  dag.add_edge
    fn (graph: dag_state, from: string, to: string) -> dag_state
    + returns a new graph with the edge added; nodes are created implicitly
    # mutation
  dag.topo_sort
    fn (graph: dag_state) -> result[list[string], string]
    + returns nodes in topological order using Kahn's algorithm
    - returns error when the graph contains a cycle, naming one cycle node
    # sorting
  dag.prune
    fn (graph: dag_state, keep: list[string]) -> dag_state
    + returns a subgraph containing only the given nodes and their transitive ancestors
    + preserves edges whose endpoints are both retained
    # pruning
