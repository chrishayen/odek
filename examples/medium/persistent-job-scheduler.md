# Requirement: "a job scheduler backed by a persistent store"

Scheduler state survives restarts by writing jobs to a pluggable key-value store. Storage and clock are injected.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.kv
    std.kv.put
      fn (store: kv_store, key: string, value: bytes) -> kv_store
      + writes value at key
      # storage
    std.kv.get
      fn (store: kv_store, key: string) -> optional[bytes]
      + reads value at key
      - returns none when key is absent
      # storage
    std.kv.delete
      fn (store: kv_store, key: string) -> kv_store
      + removes key
      # storage
    std.kv.list_prefix
      fn (store: kv_store, prefix: string) -> list[tuple[string, bytes]]
      + returns all entries whose key starts with prefix, in lexicographic order
      # storage

scheduler
  scheduler.new
    fn (store: kv_store) -> scheduler_state
    + loads any persisted jobs from store and returns a ready scheduler
    # construction
    -> std.kv.list_prefix
  scheduler.enqueue
    fn (state: scheduler_state, id: string, fire_at: i64, payload: bytes) -> scheduler_state
    + persists a new job so it survives restarts
    # registration
    -> std.kv.put
  scheduler.claim_due
    fn (state: scheduler_state) -> tuple[list[due_job], scheduler_state]
    + atomically marks all due jobs as claimed and returns them for execution
    # dispatch
    -> std.time.now_seconds
    -> std.kv.put
  scheduler.complete
    fn (state: scheduler_state, id: string) -> scheduler_state
    + removes a completed job from the store
    # lifecycle
    -> std.kv.delete
  scheduler.fail
    fn (state: scheduler_state, id: string, retry_at: i64) -> scheduler_state
    + reschedules a failed job at retry_at and clears its claim
    # lifecycle
    -> std.kv.put
