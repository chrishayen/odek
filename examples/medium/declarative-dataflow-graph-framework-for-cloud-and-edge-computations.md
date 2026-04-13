# Requirement: "a declarative framework for dataflow graphs of computations spanning cloud and edge nodes"

The library lets callers declare nodes, edges, and placements, then runs a single tick of the graph by pushing each input through its operator. No networking; execution is in-process and deterministic.

std: (all units exist)

dataflow
  dataflow.new_graph
    @ () -> graph_state
    + creates an empty graph with no nodes or edges
    # construction
  dataflow.add_source
    @ (graph: graph_state, node_id: string, placement: string) -> result[graph_state, string]
    + registers a source node that produces values
    - returns error when node_id already exists
    # topology
  dataflow.add_operator
    @ (graph: graph_state, node_id: string, placement: string) -> result[graph_state, string]
    + registers an operator node that transforms values
    - returns error when node_id already exists
    # topology
  dataflow.add_sink
    @ (graph: graph_state, node_id: string, placement: string) -> result[graph_state, string]
    + registers a sink node that collects values
    - returns error when node_id already exists
    # topology
  dataflow.connect
    @ (graph: graph_state, from_id: string, to_id: string) -> result[graph_state, string]
    + adds a directed edge between two existing nodes
    - returns error when either endpoint is unknown
    - returns error when the new edge would create a cycle
    # topology
  dataflow.topo_order
    @ (graph: graph_state) -> result[list[string], string]
    + returns node ids in topological order
    - returns error when the graph has no sources
    # scheduling
  dataflow.push
    @ (graph: graph_state, source_id: string, value: string) -> result[graph_state, string]
    + enqueues a value on a source node's output buffer
    - returns error when source_id is not a registered source
    # ingestion
  dataflow.step
    @ (graph: graph_state) -> result[graph_state, string]
    + drains buffers along the topological order, applying the identity transform at each operator
    ? operator logic is host-provided externally; this layer owns only routing
    # execution
  dataflow.drain_sink
    @ (graph: graph_state, sink_id: string) -> result[list[string], string]
    + returns and clears the values accumulated at a sink
    - returns error when sink_id is not a registered sink
    # egress
