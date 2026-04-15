# Requirement: "a distributed batch data processing framework"

Users describe a job as a dag of stages over partitioned inputs. The library schedules tasks onto workers, tracks their progress, and materializes outputs.

std
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + serializes a json value
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a json document
      # serialization
  std.net
    std.net.rpc_call
      fn (addr: string, method: string, payload: bytes) -> result[bytes, string]
      + sends a request and awaits response
      - returns error on connection failure
      # networking
    std.net.rpc_serve
      fn (addr: string, handler: rpc_handler) -> result[server_handle, string]
      + starts an rpc server
      # networking
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes a file
      # filesystem
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns the fnv-1a 64-bit hash
      # hashing

batch_framework
  batch_framework.new_job
    fn (job_name: string) -> job_graph
    + returns an empty job graph
    # construction
  batch_framework.add_stage
    fn (graph: job_graph, stage_name: string, inputs: list[string], mapper: stage_fn) -> job_graph
    + adds a stage that consumes the named inputs and produces output with the given function
    - returns a graph that will fail to validate when input names are unknown
    # job_definition
  batch_framework.validate_job
    fn (graph: job_graph) -> result[void, string]
    + checks the dag has no cycles and every input is produced or an external source
    - returns error when a cycle is present
    - returns error when an input is dangling
    # validation
  batch_framework.partition_input
    fn (input: data_source, partitions: i32) -> list[partition_ref]
    + splits an input source into partition references by key hash
    # partitioning
    -> std.hash.fnv64
  batch_framework.plan_tasks
    fn (graph: job_graph, partitions: i32) -> list[task]
    + expands every stage into one task per partition, returned in topological order
    # planning
  batch_framework.register_worker
    fn (state: coordinator_state, worker_addr: string, capacity: i32) -> coordinator_state
    + registers a worker that will be dispatched tasks
    # cluster
  batch_framework.dispatch_task
    fn (state: coordinator_state, task: task) -> result[worker_addr, string]
    + assigns a task to the worker with most spare capacity and ships the task definition
    - returns error when no workers are available
    # scheduling
    -> std.net.rpc_call
    -> std.json.encode
  batch_framework.mark_task_complete
    fn (state: coordinator_state, task_id: string, output_partitions: list[partition_ref]) -> coordinator_state
    + records completion and frees worker capacity
    # scheduling
  batch_framework.mark_task_failed
    fn (state: coordinator_state, task_id: string, reason: string) -> coordinator_state
    + records failure and queues the task for retry up to the configured limit
    # fault_tolerance
  batch_framework.run_job
    fn (state: coordinator_state, graph: job_graph) -> result[list[partition_ref], string]
    + executes the full job and returns references to its terminal outputs
    - returns error when any task exceeds retry budget
    # orchestration
  batch_framework.worker_execute_task
    fn (task: task) -> result[list[partition_ref], string]
    + loads inputs, invokes the stage function, writes outputs, and returns their refs
    - returns error when the stage function raises
    # worker
    -> std.fs.read_all
    -> std.fs.write_all
  batch_framework.worker_serve
    fn (addr: string) -> result[server_handle, string]
    + starts an rpc server that accepts task dispatches from the coordinator
    # worker
    -> std.net.rpc_serve
