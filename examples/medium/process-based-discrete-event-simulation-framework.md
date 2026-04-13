# Requirement: "a process-based discrete-event simulation framework"

A simulation environment where processes yield timeouts and events, driven by a priority queue of scheduled events.

std
  std.collections
    std.collections.min_heap_new
      @ () -> min_heap_state
      + creates an empty min-heap keyed by i64
      # collections
    std.collections.min_heap_push
      @ (heap: min_heap_state, key: i64, value: i64) -> min_heap_state
      + inserts a (key, value) pair
      # collections
    std.collections.min_heap_pop
      @ (heap: min_heap_state) -> optional[tuple[i64, i64, min_heap_state]]
      + returns the (key, value) pair with smallest key and the updated heap
      - returns none when empty
      # collections

sim
  sim.env_new
    @ () -> env_state
    + creates a simulation environment at time zero with no pending events
    # construction
    -> std.collections.min_heap_new
  sim.schedule_event
    @ (env: env_state, delay: i64, process_id: i64) -> env_state
    + schedules a wake-up for a process at now+delay
    # scheduling
    -> std.collections.min_heap_push
  sim.spawn_process
    @ (env: env_state, initial_delay: i64) -> tuple[env_state, i64]
    + registers a new process and schedules its first wake-up
    + returns (env, process_id)
    # process
  sim.timeout
    @ (env: env_state, process_id: i64, delay: i64) -> env_state
    + reschedules a process to wake after delay time units
    # process
    -> std.collections.min_heap_push
  sim.step
    @ (env: env_state) -> optional[tuple[env_state, i64, i64]]
    + advances to the next scheduled event
    + returns (env, process_id, event_time)
    - returns none when no events remain
    # scheduling
    -> std.collections.min_heap_pop
  sim.run_until
    @ (env: env_state, until: i64) -> env_state
    + repeatedly steps until the clock reaches or exceeds until
    # scheduling
  sim.now
    @ (env: env_state) -> i64
    + returns the current simulation time
    # introspection
