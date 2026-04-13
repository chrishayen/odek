# Requirement: "a job queue backed by a pluggable key/value store"

Enqueue, reserve, ack, and fail operations over a persistent store abstraction. The store is pluggable; this library encodes the queue protocol as a sequence of get/set operations the caller applies.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.new_v4
      @ () -> string
      + returns a random UUID as a string
      # identifiers
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

job_queue
  job_queue.new
    @ (queue_name: string) -> queue_state
    + creates a named queue handle
    # construction
  job_queue.enqueue
    @ (state: queue_state, payload: map[string, string]) -> tuple[string, queue_state]
    + assigns a job id and records the job as pending
    # enqueue
    -> std.uuid.new_v4
    -> std.time.now_millis
    -> std.json.encode_object
  job_queue.reserve
    @ (state: queue_state, worker_id: string, lease_ms: i64) -> tuple[optional[reserved_job], queue_state]
    + returns the next pending job and marks it reserved until the lease expires
    - returns none when no job is available
    # reservation
    -> std.time.now_millis
  job_queue.ack
    @ (state: queue_state, job_id: string) -> queue_state
    + marks a reserved job as completed
    # completion
  job_queue.fail
    @ (state: queue_state, job_id: string, reason: string) -> queue_state
    + increments attempt counter and requeues or dead-letters based on max attempts
    # failure
  job_queue.expire_leases
    @ (state: queue_state) -> queue_state
    + returns expired reservations to the pending set
    # recovery
    -> std.time.now_millis
  job_queue.decode_job
    @ (raw: string) -> result[map[string, string], string]
    + decodes a stored job payload
    - returns error when the stored form is malformed
    # serialization
    -> std.json.parse_object
