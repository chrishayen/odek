# Requirement: "a job queue engine backed by a key-value store with retries, delays, and priorities"

Stores jobs in a pluggable kv backend and exposes enqueue/dequeue/ack. Project layer owns the job lifecycle; std provides time, ids, and json.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_ulid
      fn () -> string
      + returns a newly generated ulid string
      # identifiers
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

job_queue
  job_queue.new
    fn (backend: kv_backend) -> queue_state
    + wraps a kv backend with queue metadata
    # construction
  job_queue.enqueue
    fn (q: queue_state, task_type: string, payload: map[string, string], priority: i32, delay_ms: i64) -> result[string, string]
    + stores a job and returns its id
    + jobs with delay_ms > 0 are invisible until ready_at is reached
    # enqueue
    -> std.id.new_ulid
    -> std.time.now_millis
    -> std.json.encode_object
  job_queue.fetch_ready
    fn (q: queue_state) -> result[optional[job], string]
    + returns the highest-priority job whose ready_at has elapsed
    + marks the fetched job as in-flight with a lease deadline
    - returns none when nothing is ready
    # dequeue
    -> std.time.now_millis
    -> std.json.parse_object
  job_queue.ack
    fn (q: queue_state, job_id: string) -> result[queue_state, string]
    + removes the job from in-flight storage
    - returns error when the job is not in-flight
    # completion
  job_queue.nack
    fn (q: queue_state, job_id: string, error_msg: string) -> result[queue_state, string]
    + increments the job's attempt count and either requeues with backoff or moves it to the failed set
    # retry
    -> std.time.now_millis
  job_queue.reclaim_expired
    fn (q: queue_state) -> result[queue_state, string]
    + returns expired in-flight jobs to the ready queue
    # leases
    -> std.time.now_millis
  job_queue.stats
    fn (q: queue_state) -> queue_stats
    + returns counts of ready, in-flight, delayed, and failed jobs
    # observability
