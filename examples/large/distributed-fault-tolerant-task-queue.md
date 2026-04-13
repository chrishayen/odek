# Requirement: "a distributed fault-tolerant task queue"

Producers enqueue named tasks with payloads. Workers lease tasks with a visibility timeout; failed or lost leases are retried with backoff until a retry ceiling, then parked in a dead-letter queue.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + current unix time in milliseconds
      # time
  std.id
    std.id.random_128
      @ () -> string
      + returns a random 128-bit id in a stable textual form
      # identifiers
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + serializes json
      # serialization
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses json
      # serialization
  std.store
    std.store.tx_begin
      @ (db: store_handle) -> result[tx, string]
      + begins a transaction
      # storage
    std.store.tx_commit
      @ (t: tx) -> result[void, string]
      + commits a transaction
      # storage
    std.store.tx_rollback
      @ (t: tx) -> void
      + discards a transaction
      # storage
    std.store.put
      @ (t: tx, key: bytes, value: bytes) -> result[void, string]
      + writes a key/value
      # storage
    std.store.get
      @ (t: tx, key: bytes) -> result[optional[bytes], string]
      + reads a key
      # storage
    std.store.range
      @ (t: tx, prefix: bytes) -> list[tuple[bytes, bytes]]
      + returns all key/value pairs with the given prefix
      # storage

task_queue
  task_queue.open
    @ (db: store_handle) -> queue_state
    + returns a queue handle backed by the given store
    # construction
  task_queue.enqueue
    @ (state: queue_state, task_type: string, payload: bytes) -> result[string, string]
    + persists a new task with status pending and returns its id
    + records created_at and retry_count=0
    # enqueue
    -> std.id.random_128
    -> std.time.now_millis
    -> std.store.tx_begin
    -> std.store.put
    -> std.store.tx_commit
  task_queue.lease
    @ (state: queue_state, worker_id: string, visibility_ms: i64) -> result[optional[leased_task], string]
    + atomically claims the oldest pending or expired task and marks it in-flight with a new lease deadline
    + returns none when no task is available
    # lease
    -> std.time.now_millis
    -> std.store.tx_begin
    -> std.store.range
    -> std.store.put
    -> std.store.tx_commit
  task_queue.ack
    @ (state: queue_state, task_id: string, worker_id: string) -> result[void, string]
    + deletes the task when the lease is still owned by worker_id
    - returns error when the lease has expired or was stolen
    # completion
    -> std.store.tx_begin
    -> std.store.get
    -> std.store.put
    -> std.store.tx_commit
  task_queue.nack
    @ (state: queue_state, task_id: string, worker_id: string, reason: string) -> result[void, string]
    + increments retry_count and reschedules with exponential backoff
    + when retry_count exceeds the ceiling the task moves to the dead-letter queue
    # retry
    -> std.time.now_millis
  task_queue.compute_backoff_ms
    @ (retry_count: i32) -> i64
    + returns an exponential backoff with jitter for the given retry count
    # retry
  task_queue.reap_expired
    @ (state: queue_state) -> result[i32, string]
    + finds in-flight tasks whose lease deadline has passed and resets them to pending
    + returns the number of tasks reaped
    # fault_tolerance
    -> std.time.now_millis
  task_queue.list_dead_letters
    @ (state: queue_state) -> result[list[task_summary], string]
    + returns all tasks parked in the dead-letter queue
    # inspection
    -> std.store.range
  task_queue.requeue_dead_letter
    @ (state: queue_state, task_id: string) -> result[void, string]
    + moves a dead-lettered task back into the pending queue with retry_count reset
    - returns error when the task id is not in the dead-letter queue
    # recovery
