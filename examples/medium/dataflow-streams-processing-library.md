# Requirement: "a dataflow streams processing library"

Named blocks are wired together into a graph; each block transforms messages and forwards them along outgoing edges. No UI, just the graph and its engine.

std: (all units exist)

streams
  streams.new_graph
    @ () -> graph_state
    + returns an empty graph with no blocks or edges
    # construction
  streams.add_block
    @ (g: graph_state, id: string, kind: string, config: map[string, string]) -> result[graph_state, string]
    + registers a block under id with a transform kind and configuration
    - returns error when id is already used
    - returns error when kind is not a known block type
    # topology
  streams.connect
    @ (g: graph_state, from_id: string, to_id: string) -> result[graph_state, string]
    + adds a directed edge from one block to another
    - returns error when either id is unknown
    - returns error when the edge would create a cycle
    # topology
  streams.send
    @ (g: graph_state, target_id: string, message: bytes) -> result[graph_state, string]
    + enqueues a message on the input of the target block
    - returns error when target_id is unknown
    # ingestion
  streams.tick
    @ (g: graph_state) -> graph_state
    + drains every block's input queue once, applying its transform and forwarding results downstream
    # execution
  streams.drain
    @ (g: graph_state, sink_id: string) -> tuple[graph_state, list[bytes]]
    + removes and returns every message currently buffered at the sink block
    # output
