# Requirement: "a graph library for analyzing complex networks"

Directed and undirected graph data structure with common analysis algorithms.

std: (all units exist)

graph
  graph.new
    @ (directed: bool) -> graph_state
    + creates an empty graph
    # construction
  graph.add_node
    @ (g: graph_state, id: string) -> graph_state
    + inserts a node; no-op if already present
    # mutation
  graph.add_edge
    @ (g: graph_state, src: string, dst: string, weight: f64) -> graph_state
    + adds an edge, inserting either endpoint if missing
    + adds the reverse edge automatically when the graph is undirected
    # mutation
  graph.neighbors
    @ (g: graph_state, id: string) -> list[string]
    + returns outgoing neighbors for id
    - returns empty list when id is absent
    # traversal
  graph.bfs
    @ (g: graph_state, start: string) -> list[string]
    + returns nodes in breadth-first order from start
    # traversal
    -> graph.neighbors
  graph.dfs
    @ (g: graph_state, start: string) -> list[string]
    + returns nodes in depth-first pre-order from start
    # traversal
    -> graph.neighbors
  graph.shortest_path
    @ (g: graph_state, src: string, dst: string) -> optional[list[string]]
    + returns the minimum-weight path using Dijkstra's algorithm
    - returns none when dst is unreachable
    ? all edge weights are assumed non-negative
    # pathfinding
    -> graph.neighbors
  graph.connected_components
    @ (g: graph_state) -> list[list[string]]
    + returns groups of nodes mutually reachable along undirected edges
    # components
    -> graph.neighbors
  graph.degree_centrality
    @ (g: graph_state) -> map[string, f64]
    + returns each node's degree normalized by (n - 1)
    # centrality
  graph.pagerank
    @ (g: graph_state, damping: f64, iterations: i32) -> map[string, f64]
    + runs the power iteration for the given number of steps
    + returned values sum to 1.0
    # centrality
  graph.clustering_coefficient
    @ (g: graph_state, id: string) -> f64
    + returns the fraction of possible edges among id's neighbors that exist
    - returns 0 when id has fewer than two neighbors
    # clustering
    -> graph.neighbors
  graph.topological_sort
    @ (g: graph_state) -> result[list[string], string]
    + returns nodes in a valid topological order
    - returns error when the graph contains a cycle or is undirected
    # ordering
    -> graph.neighbors
