# Requirement: "a type-safe, persistent, embedded task queue and background job runner backed by an embedded SQL database"

Jobs are enqueued with a typed payload, persisted to an embedded SQL store, and a worker pulls ready jobs off the queue, runs them, and records the outcome.

std
  std.sql
    std.sql.open
      fn (path: string) -> result[sql_conn, string]
      + opens or creates an embedded SQL database file
      - returns error when the path cannot be opened
      # storage
    std.sql.exec
      fn (conn: sql_conn, query: string, args: list[string]) -> result[i64, string]
      + runs a statement and returns rows affected
      - returns error on invalid SQL
      # storage
    std.sql.query_rows
      fn (conn: sql_conn, query: string, args: list[string]) -> result[list[map[string,string]], string]
      + runs a query and returns rows as column-name maps
      - returns error on invalid SQL
      # storage
    std.sql.close
      fn (conn: sql_conn) -> void
      + releases the underlying database handle
      # storage
  std.time
    std.time.now_millis
      fn () -> i64
      + returns the current unix time in milliseconds
      # time
  std.random
    std.random.uuid_v4
      fn () -> string
      + returns a random 128-bit identifier as a hex string
      # identifiers

task_queue
  task_queue.open
    fn (path: string) -> result[queue_state, string]
    + opens a queue backed by a SQL file and creates the jobs table if missing
    - returns error when the underlying store cannot be opened
    # construction
    -> std.sql.open
    -> std.sql.exec
  task_queue.register
    fn (state: queue_state, kind: string, schema_keys: list[string]) -> queue_state
    + records that a job kind is valid and what payload keys it expects
    ? unknown kinds are rejected at enqueue time to keep the queue type-safe
    # registration
  task_queue.enqueue
    fn (state: queue_state, kind: string, payload: map[string,string], run_at_ms: i64) -> result[string, string]
    + inserts a pending job and returns its id
    - returns error when kind is not registered
    - returns error when payload is missing a required key
    # enqueue
    -> std.random.uuid_v4
    -> std.sql.exec
  task_queue.claim_next
    fn (state: queue_state, worker_id: string) -> result[optional[claimed_job], string]
    + atomically marks the next ready job as running and returns it
    + returns none when no job is ready to run
    ? "ready" means pending with run_at_ms <= now
    # claim
    -> std.time.now_millis
    -> std.sql.exec
    -> std.sql.query_rows
  task_queue.mark_done
    fn (state: queue_state, job_id: string) -> result[void, string]
    + marks a claimed job as completed
    - returns error when the job id does not exist
    # completion
    -> std.sql.exec
  task_queue.mark_failed
    fn (state: queue_state, job_id: string, reason: string, retry_in_ms: i64) -> result[void, string]
    + moves a job back to pending with a delayed run_at when retry_in_ms > 0
    + marks the job as dead-lettered when retry_in_ms <= 0
    # failure
    -> std.time.now_millis
    -> std.sql.exec
  task_queue.run_worker_once
    fn (state: queue_state, worker_id: string, handler: fn(string, map[string,string]) -> result[void,string]) -> result[bool, string]
    + claims one ready job, invokes the handler, and records success or failure
    + returns false when no job was ready
    # worker_loop
    -> std.time.now_millis
  task_queue.stats
    fn (state: queue_state) -> result[map[string,i64], string]
    + returns counts by status (pending, running, done, dead)
    # introspection
    -> std.sql.query_rows
  task_queue.close
    fn (state: queue_state) -> void
    + flushes and releases the underlying connection
    # teardown
    -> std.sql.close
