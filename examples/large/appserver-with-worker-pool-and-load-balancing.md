# Requirement: "an application server with worker pool, load balancing, and process supervision"

A worker-pool server that dispatches work items across supervised worker processes with configurable load-balancing strategies and health-based restart.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.new_v4
      fn () -> string
      + returns a random UUID as a string
      # identifiers

appserver
  appserver.new_pool
    fn (min_workers: i32, max_workers: i32) -> pool_state
    + creates a pool with size bounds and no workers yet
    # construction
  appserver.spawn_worker
    fn (state: pool_state, command: string) -> tuple[string, pool_state]
    + records a new worker with the given command and returns its id
    - returns unchanged state when the pool is at max_workers
    # supervision
    -> std.uuid.new_v4
    -> std.time.now_millis
  appserver.mark_worker_ready
    fn (state: pool_state, worker_id: string) -> pool_state
    + transitions a worker from starting to ready
    # supervision
    -> std.time.now_millis
  appserver.mark_worker_exited
    fn (state: pool_state, worker_id: string, exit_code: i32) -> pool_state
    + marks a worker exited and schedules restart if abnormal
    # supervision
    -> std.time.now_millis
  appserver.restart_dead
    fn (state: pool_state) -> tuple[list[string], pool_state]
    + restarts workers that exited abnormally, returning new ids
    # supervision
    -> std.uuid.new_v4
  appserver.set_strategy
    fn (state: pool_state, strategy: string) -> result[pool_state, string]
    + selects load-balancing strategy: "round_robin", "least_busy", or "random"
    - returns error on unknown strategy name
    # load_balancing
  appserver.submit
    fn (state: pool_state, work_item_id: string) -> result[tuple[string, pool_state], string]
    + returns the chosen worker id and updated state
    - returns error when no workers are ready
    # dispatch
    -> std.time.now_millis
  appserver.complete
    fn (state: pool_state, worker_id: string, work_item_id: string) -> pool_state
    + records work completion and frees the worker
    # dispatch
    -> std.time.now_millis
  appserver.health_check
    fn (state: pool_state, worker_id: string, healthy: bool) -> pool_state
    + records a health-check result; unhealthy workers are drained
    # supervision
    -> std.time.now_millis
  appserver.drain
    fn (state: pool_state) -> pool_state
    + marks all workers drain-pending so no new work is dispatched
    # lifecycle
  appserver.stats
    fn (state: pool_state) -> pool_stats
    + returns counts of ready, busy, and exited workers
    # observability
