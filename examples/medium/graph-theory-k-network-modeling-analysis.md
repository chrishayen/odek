# Requirement: "a graph theory modeling and analysis library"

Build directed or undirected graphs with labeled nodes and weighted edges, then run standard analyses: shortest path, connected components, and centrality.

std: (all units exist)

graph_lib
  graph_lib.new
    @ (directed: bool) -> graph_state
    + creates an empty graph with the given directedness
    # construction
  graph_lib.add_node
    @ (state: graph_state, id: string) -> graph_state
    + adds a node with the given id
    - returns unchanged state when the node already exists
    # mutation
  graph_lib.add_edge
    @ (state: graph_state, from: string, to: string, weight: f64) -> result[graph_state, string]
    + adds a weighted edge between two nodes
    - returns error when either endpoint is unknown
    # mutation
  graph_lib.neighbors
    @ (state: graph_state, id: string) -> list[string]
    + returns the ids of nodes adjacent to id
    # query
  graph_lib.shortest_path
    @ (state: graph_state, from: string, to: string) -> optional[list[string]]
    + returns the minimum-weight path as a list of node ids
    - returns none when no path exists
    ? uses Dijkstra; requires non-negative edge weights
    # analysis
  graph_lib.connected_components
    @ (state: graph_state) -> list[list[string]]
    + returns the groups of mutually reachable nodes
    + for directed graphs, returns strongly connected components
    # analysis
  graph_lib.degree_centrality
    @ (state: graph_state) -> map[string, f64]
    + returns each node's normalized degree centrality
    # analysis
  graph_lib.betweenness_centrality
    @ (state: graph_state) -> map[string, f64]
    + returns each node's betweenness centrality score
    # analysis
