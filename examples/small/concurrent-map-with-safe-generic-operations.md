# Requirement: "a concurrent map with safe generic operations"

A thread-safe map backed by a mutex from std.

std
  std.sync
    std.sync.new_mutex
      fn () -> mutex_handle
      + returns a new unlocked mutex
      # concurrency
    std.sync.lock
      fn (m: mutex_handle) -> void
      + blocks until the mutex is acquired
      # concurrency
    std.sync.unlock
      fn (m: mutex_handle) -> void
      + releases the mutex
      # concurrency

concurrent_map
  concurrent_map.new
    fn () -> cmap_state
    + returns an empty concurrent map
    # construction
    -> std.sync.new_mutex
  concurrent_map.put
    fn (m: cmap_state, key: string, value: bytes) -> cmap_state
    + inserts or updates a key atomically
    # mutation
    -> std.sync.lock
    -> std.sync.unlock
  concurrent_map.get
    fn (m: cmap_state, key: string) -> optional[bytes]
    + returns the value for a key under a lock
    - returns none when the key is absent
    # reads
    -> std.sync.lock
    -> std.sync.unlock
  concurrent_map.delete
    fn (m: cmap_state, key: string) -> cmap_state
    + removes a key atomically
    # mutation
    -> std.sync.lock
    -> std.sync.unlock
  concurrent_map.compute_if_absent
    fn (m: cmap_state, key: string, default_value: bytes) -> tuple[bytes, cmap_state]
    + returns the existing value or inserts and returns the default, atomically
    # mutation
    -> std.sync.lock
    -> std.sync.unlock
