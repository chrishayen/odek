# Requirement: "a parallel computing library for analytic workloads"

A task graph executor: users build a DAG of typed computations, the library partitions it and runs tasks across workers.

std
  std.concurrency
    std.concurrency.spawn_workers
      fn (count: i32) -> worker_pool
      + starts a worker pool of the given size
      # concurrency
    std.concurrency.submit
      fn (pool: worker_pool, task: work_task) -> task_handle
      + submits a task and returns a handle for awaiting its result
      # concurrency
    std.concurrency.await
      fn (handle: task_handle) -> result[work_value, string]
      + blocks until the task completes and returns its result
      - returns error when the task panicked or was cancelled
      # concurrency
    std.concurrency.shutdown
      fn (pool: worker_pool) -> void
      + drains the pool and releases its workers
      # concurrency
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns a 64-bit FNV-1a hash of the input
      # hashing

dataflow
  dataflow.new_graph
    fn () -> graph
    + creates an empty task graph
    # construction
  dataflow.add_source
    fn (g: graph, name: string, values: list[work_value]) -> graph
    + adds a source node producing the given values
    # graph
  dataflow.add_map
    fn (g: graph, name: string, input: string, fn: map_fn) -> graph
    + adds a map node that transforms each value from input
    - returns unchanged graph (tagged) when input is unknown
    # graph
  dataflow.add_filter
    fn (g: graph, name: string, input: string, fn: filter_fn) -> graph
    + adds a filter node keeping values where fn returns true
    # graph
  dataflow.add_reduce
    fn (g: graph, name: string, input: string, init: work_value, fn: reduce_fn) -> graph
    + adds a reduce node that folds input values into a single result
    # graph
  dataflow.add_join
    fn (g: graph, name: string, left: string, right: string, key_fn: key_fn) -> graph
    + adds a join node keyed by key_fn; uses hashing for partitioning
    # graph
    -> std.hash.fnv64
  dataflow.validate
    fn (g: graph) -> result[void, string]
    + checks for cycles and missing dependencies
    - returns error listing the cycle nodes
    - returns error naming a missing input
    # validation
  dataflow.plan
    fn (g: graph) -> execution_plan
    + topologically orders nodes into stages executable in parallel
    # planning
  dataflow.execute
    fn (plan: execution_plan, worker_count: i32) -> result[map[string, list[work_value]], string]
    + runs the plan and returns per-node output values
    - returns error when a task fails
    # execution
    -> std.concurrency.spawn_workers
    -> std.concurrency.submit
    -> std.concurrency.await
    -> std.concurrency.shutdown
  dataflow.get_output
    fn (results: map[string, list[work_value]], name: string) -> result[list[work_value], string]
    + returns the output of a named node
    - returns error when the node is missing
    # access
