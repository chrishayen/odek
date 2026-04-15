# Requirement: "a concurrent-safe generic data structure library"

Thread-safe variants of common collections.

std
  std.sync
    std.sync.new_mutex
      fn () -> mutex_state
      + returns an unlocked mutex
      # concurrency
    std.sync.with_lock
      fn (m: mutex_state, fn: lock_fn) -> void
      + acquires the mutex, runs fn, then releases
      # concurrency

concurrent_collections
  concurrent_collections.new_map
    fn () -> safe_map_state
    + creates an empty concurrent key-value map
    # construction
    -> std.sync.new_mutex
  concurrent_collections.map_set
    fn (m: safe_map_state, key: string, value: bytes) -> safe_map_state
    + inserts or overwrites a key
    # mutation
    -> std.sync.with_lock
  concurrent_collections.map_get
    fn (m: safe_map_state, key: string) -> optional[bytes]
    + returns the stored value when present
    - returns none when the key is absent
    # lookup
    -> std.sync.with_lock
  concurrent_collections.new_queue
    fn () -> safe_queue_state
    + creates an empty concurrent fifo queue
    # construction
    -> std.sync.new_mutex
  concurrent_collections.queue_push
    fn (q: safe_queue_state, item: bytes) -> safe_queue_state
    + appends an item to the tail
    # mutation
    -> std.sync.with_lock
  concurrent_collections.queue_pop
    fn (q: safe_queue_state) -> optional[bytes]
    + removes and returns the head item
    - returns none when the queue is empty
    # mutation
    -> std.sync.with_lock
  concurrent_collections.new_set
    fn () -> safe_set_state
    + creates an empty concurrent set of strings
    # construction
    -> std.sync.new_mutex
  concurrent_collections.set_add
    fn (s: safe_set_state, item: string) -> safe_set_state
    + inserts an item if not already present
    # mutation
    -> std.sync.with_lock
  concurrent_collections.set_contains
    fn (s: safe_set_state, item: string) -> bool
    + returns true when the item is in the set
    # lookup
    -> std.sync.with_lock
