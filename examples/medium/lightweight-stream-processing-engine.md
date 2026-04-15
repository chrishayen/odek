# Requirement: "a library for a lightweight stream processing engine"

Builds a small dataflow graph from sources through transforms to sinks and runs tuples through it.

std: (all units exist)

stream_engine
  stream_engine.new_topology
    fn () -> topology
    + creates an empty dataflow graph
    # construction
  stream_engine.add_source
    fn (topo: topology, name: string, fetch: tuple_fetcher) -> result[topology, string]
    + registers a named source that produces tuples on demand
    - returns error when name is already used
    # topology_building
  stream_engine.add_map
    fn (topo: topology, name: string, input: string, fn: tuple_transform) -> result[topology, string]
    + registers a map operator that transforms each tuple from input
    - returns error when input does not refer to an existing node
    # topology_building
  stream_engine.add_filter
    fn (topo: topology, name: string, input: string, pred: tuple_predicate) -> result[topology, string]
    + registers a filter operator that keeps only tuples matching the predicate
    - returns error when input does not refer to an existing node
    # topology_building
  stream_engine.add_window
    fn (topo: topology, name: string, input: string, size: i32) -> result[topology, string]
    + registers a tumbling window that emits a list of tuples once it has size elements
    - returns error when size is not positive
    # topology_building
  stream_engine.add_sink
    fn (topo: topology, name: string, input: string, fn: tuple_sink) -> result[topology, string]
    + registers a terminal sink fed by input
    - returns error when input does not refer to an existing node
    # topology_building
  stream_engine.run_once
    fn (topo: topology) -> result[i32, string]
    + pulls one tuple from every source and propagates it through the graph in topological order
    + returns the number of tuples delivered to sinks
    - returns error on the first operator failure
    # execution
