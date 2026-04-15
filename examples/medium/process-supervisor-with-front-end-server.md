# Requirement: "a process supervisor that starts, monitors, and restarts worker processes behind a front-end server"

Tracks a pool of worker processes, restarts crashed workers, and exposes a routing decision for an incoming connection.

std
  std.process
    std.process.spawn
      fn (command: string, args: list[string]) -> result[i32, string]
      + starts the command with the given args and returns its pid
      - returns error when the command cannot be launched
      # process
    std.process.is_alive
      fn (pid: i32) -> bool
      + returns true when the pid refers to a running process
      # process
    std.process.kill
      fn (pid: i32) -> result[void, string]
      + sends a terminate signal to the pid
      - returns error when the pid does not exist
      # process

process_supervisor
  process_supervisor.new_pool
    fn (command: string, args: list[string], target_size: i32) -> pool_state
    + returns a pool configured with a worker command and target worker count
    # construction
  process_supervisor.ensure_workers
    fn (pool: pool_state) -> result[pool_state, string]
    + spawns new workers until the number of live workers matches the target
    - returns error when a spawn attempt fails
    # supervision
    -> std.process.spawn
  process_supervisor.reap_dead
    fn (pool: pool_state) -> pool_state
    + removes pids from the pool that are no longer alive
    # supervision
    -> std.process.is_alive
  process_supervisor.pick_worker
    fn (pool: pool_state) -> optional[i32]
    + returns the next worker pid in round-robin order
    - returns none when the pool is empty
    # routing
  process_supervisor.shutdown
    fn (pool: pool_state) -> result[pool_state, string]
    + kills every worker in the pool and returns an empty pool
    - returns error when any kill attempt fails
    # lifecycle
    -> std.process.kill
