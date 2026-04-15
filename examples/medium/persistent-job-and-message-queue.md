# Requirement: "a persistent job and message queue"

Enqueues work for later processing with durability across process restarts, retries, and delayed execution.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.kv
    std.kv.open
      fn (path: string) -> result[kv_store, string]
      + opens or creates a durable key-value store rooted at path
      - returns error when path is unusable
      # storage
    std.kv.put
      fn (store: kv_store, key: bytes, value: bytes) -> result[void, string]
      + writes a key, fsyncing before returning
      # storage
    std.kv.get
      fn (store: kv_store, key: bytes) -> result[optional[bytes], string]
      + returns the stored value or none
      # storage
    std.kv.delete
      fn (store: kv_store, key: bytes) -> result[void, string]
      + removes the key
      # storage
    std.kv.scan_prefix
      fn (store: kv_store, prefix: bytes) -> result[list[tuple[bytes, bytes]], string]
      + returns key-value pairs whose key starts with prefix, in key order
      # storage

job_queue
  job_queue.open
    fn (path: string) -> result[queue_state, string]
    + opens a queue backed by durable storage at path
    - returns error when the backing store cannot be opened
    # construction
    -> std.kv.open
  job_queue.enqueue
    fn (state: queue_state, topic: string, payload: bytes, run_at_ms: i64) -> result[string, string]
    + inserts a job on the topic, scheduled for run_at_ms
    + returns a unique job id
    ? run_at_ms <= now means "run immediately"
    # enqueue
    -> std.time.now_millis
    -> std.kv.put
  job_queue.claim
    fn (state: queue_state, topic: string, worker_id: string, lease_ms: i32) -> result[optional[job], string]
    + returns the next ready job on the topic and marks it leased to worker_id
    + returns none when no ready jobs are present
    # claim
    -> std.time.now_millis
    -> std.kv.scan_prefix
    -> std.kv.put
  job_queue.complete
    fn (state: queue_state, job_id: string) -> result[void, string]
    + removes the job from the queue
    - returns error when the job id is unknown
    # completion
    -> std.kv.delete
  job_queue.fail
    fn (state: queue_state, job_id: string, retry_delay_ms: i32) -> result[void, string]
    + returns the job to the queue for a retry after the delay
    - returns error when the job has exceeded its retry budget
    # retry
    -> std.kv.put
  job_queue.pending_count
    fn (state: queue_state, topic: string) -> result[i64, string]
    + returns the number of jobs currently pending on the topic
    # metrics
    -> std.kv.scan_prefix
