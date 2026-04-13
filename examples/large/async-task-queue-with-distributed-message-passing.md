# Requirement: "an asynchronous task queue backed by distributed message passing"

Workers pull encoded tasks from a pluggable broker and result backend; the project wires task registration, enqueue, and worker loop.

std
  std.json
    std.json.encode_object
      @ (obj: map[string,string]) -> string
      + encodes a flat string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string map
      - returns error on non-object root
      # serialization
  std.uuid
    std.uuid.new_v4
      @ () -> string
      + returns a random UUIDv4 string
      # identifiers
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

taskq
  taskq.queue_new
    @ (broker: broker_handle, backend: backend_handle) -> queue_state
    + creates a queue with broker and result backend attached
    # construction
  taskq.register
    @ (queue: queue_state, name: string, handler: task_handler) -> queue_state
    + associates a handler with a task name
    # registration
  taskq.enqueue
    @ (queue: queue_state, name: string, args: map[string,string]) -> result[string, string]
    + publishes a task message and returns the generated task id
    - returns error when the task name is not registered
    # enqueue
    -> std.uuid.new_v4
    -> std.json.encode_object
    -> std.time.now_seconds
  taskq.fetch_next
    @ (queue: queue_state) -> result[optional[task_envelope], string]
    + pops the next task from the broker, blocking briefly if empty
    # fetch
    -> std.json.parse_object
  taskq.execute
    @ (queue: queue_state, envelope: task_envelope) -> result[string, string]
    + looks up the handler and runs it, returning the serialized result
    - returns error when the handler is missing
    # execution
  taskq.store_result
    @ (queue: queue_state, task_id: string, value: string) -> result[void, string]
    + writes the result to the backend keyed by task id
    # results
  taskq.get_result
    @ (queue: queue_state, task_id: string) -> result[optional[string], string]
    + reads a completed result, if present
    # results
  taskq.worker_run
    @ (queue: queue_state) -> result[void, string]
    + loops fetching, executing, and storing results until shutdown
    - returns error when the broker connection is lost
    # worker
