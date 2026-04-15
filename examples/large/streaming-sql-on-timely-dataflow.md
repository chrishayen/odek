# Requirement: "a streaming SQL database built on timely dataflow"

A streaming SQL engine that compiles queries into dataflow graphs and maintains incremental results as sources emit changes. Large system with a parser, planner, dataflow runtime, and source/sink plumbing.

std
  std.hashing
    std.hashing.hash_bytes
      fn (data: bytes) -> u64
      + returns a 64-bit hash suitable for partitioning and keying
      # hashing
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

streamsql
  streamsql.parse_query
    fn (sql: string) -> result[query_ast, string]
    + parses a SELECT/JOIN/GROUP BY query into an AST
    - returns error on syntax errors with position info
    # parsing
  streamsql.plan_query
    fn (ast: query_ast, catalog: catalog_state) -> result[dataflow_plan, string]
    + produces a logical plan of dataflow operators (map, filter, join, reduce)
    - returns error when a referenced source does not exist in the catalog
    # planning
  streamsql.compile_plan
    fn (plan: dataflow_plan) -> dataflow_graph
    + materializes operators into a graph wired to input/output channels
    # compilation
  streamsql.register_source
    fn (catalog: catalog_state, name: string, schema: list[string]) -> catalog_state
    + adds a named input source with its column schema
    # catalog
  streamsql.new_runtime
    fn (graph: dataflow_graph, workers: i32) -> runtime_state
    + creates a multi-worker runtime ready to accept input batches
    ? timestamps flow alongside records so operators can advance frontiers
    # runtime
  streamsql.push_batch
    fn (runtime: runtime_state, source: string, rows: list[list[string]], timestamp: i64) -> result[runtime_state, string]
    + ingests a batch of rows at the given logical timestamp
    - returns error when source name is unknown
    # ingestion
    -> std.hashing.hash_bytes
  streamsql.advance_frontier
    fn (runtime: runtime_state, source: string, timestamp: i64) -> runtime_state
    + signals that no more records for source will arrive before timestamp
    ? downstream reducers emit final results when all inputs pass the frontier
    # frontier
  streamsql.drain_output
    fn (runtime: runtime_state, sink: string) -> tuple[list[output_change], runtime_state]
    + returns pending +/- change records for a named sink and clears them
    # output
  streamsql.snapshot_results
    fn (runtime: runtime_state, sink: string) -> list[list[string]]
    + returns the current materialized rows for a sink as of the latest frontier
    # snapshot
  streamsql.shutdown_runtime
    fn (runtime: runtime_state) -> void
    + tears down workers and releases graph resources
    # lifecycle
    -> std.time.now_millis
