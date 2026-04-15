# Requirement: "a basic graph algorithms library"

A small directed-graph data structure and the standard traversal and shortest-path algorithms over it.

std: (all units exist)

graph
  graph.new
    fn () -> graph_state
    + returns an empty directed graph
    # construction
  graph.add_node
    fn (g: graph_state, id: string) -> graph_state
    + adds a node with the given id
    ? adding an existing id is a no-op
    # construction
  graph.add_edge
    fn (g: graph_state, from: string, to: string, weight: f64) -> graph_state
    + adds a weighted directed edge between two nodes
    # construction
  graph.neighbors
    fn (g: graph_state, id: string) -> list[string]
    + returns the ids of nodes reachable via a single outgoing edge
    + returns an empty list when the node has no outgoing edges
    # query
  graph.bfs
    fn (g: graph_state, start: string) -> list[string]
    + returns node ids in breadth-first order starting from start
    - returns an empty list when start is not in the graph
    # traversal
  graph.dfs
    fn (g: graph_state, start: string) -> list[string]
    + returns node ids in depth-first order starting from start
    # traversal
  graph.shortest_path
    fn (g: graph_state, from: string, to: string) -> optional[list[string]]
    + returns the sequence of nodes forming the shortest weighted path
    + returns none when no path exists
    ? edge weights must be non-negative
    # shortest_path
  graph.topological_sort
    fn (g: graph_state) -> result[list[string], string]
    + returns nodes in a topological order
    - returns error when the graph contains a cycle
    # ordering
  graph.connected_components
    fn (g: graph_state) -> list[list[string]]
    + returns the weakly connected components of the graph
    # components
