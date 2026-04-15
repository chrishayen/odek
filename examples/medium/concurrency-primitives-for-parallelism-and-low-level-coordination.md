# Requirement: "concurrency primitives for parallelism and low-level coordination"

A small toolkit of concurrency primitives: a multi-producer/multi-consumer channel, an atomic counter, and a wait group. Threading and atomics are host-level; the project exposes composed primitives.

std
  std.atomic
    std.atomic.load_i64
      fn (addr: atomic_i64) -> i64
      + returns the current value with acquire ordering
      # atomic
    std.atomic.store_i64
      fn (addr: atomic_i64, value: i64) -> void
      + stores the value with release ordering
      # atomic
    std.atomic.cas_i64
      fn (addr: atomic_i64, expected: i64, new_value: i64) -> bool
      + returns true when the compare-and-swap succeeded
      - returns false when the current value did not match expected
      # atomic
  std.sync
    std.sync.park
      fn () -> void
      + parks the current thread until unparked
      # threading
    std.sync.unpark
      fn (thread: thread_handle) -> void
      + wakes the parked thread
      # threading

concurrency
  concurrency.channel_new
    fn (capacity: i32) -> channel_state
    + creates a bounded channel; capacity 0 means unbuffered
    # construction
  concurrency.channel_send
    fn (ch: channel_state, value: bytes) -> result[void, string]
    + enqueues the value, blocking when full
    - returns error when the channel has been closed
    # messaging
    -> std.atomic.cas_i64
    -> std.sync.park
  concurrency.channel_recv
    fn (ch: channel_state) -> result[bytes, string]
    + dequeues a value, blocking when empty
    - returns error when the channel is closed and drained
    # messaging
    -> std.atomic.cas_i64
    -> std.sync.unpark
  concurrency.channel_close
    fn (ch: channel_state) -> void
    + marks the channel closed and wakes all waiters
    # lifecycle
    -> std.sync.unpark
  concurrency.counter_new
    fn (initial: i64) -> counter_state
    + creates an atomic counter seeded to initial
    # construction
  concurrency.counter_add
    fn (c: counter_state, delta: i64) -> i64
    + atomically adds delta and returns the new value
    # counter
    -> std.atomic.cas_i64
  concurrency.waitgroup_new
    fn () -> waitgroup_state
    + creates a wait group with zero pending tasks
    # construction
  concurrency.waitgroup_add
    fn (wg: waitgroup_state, delta: i32) -> void
    + increments the pending count
    # sync
    -> std.atomic.store_i64
  concurrency.waitgroup_done
    fn (wg: waitgroup_state) -> void
    + decrements pending and wakes waiters when zero
    # sync
    -> std.atomic.cas_i64
    -> std.sync.unpark
  concurrency.waitgroup_wait
    fn (wg: waitgroup_state) -> void
    + blocks until pending reaches zero
    # sync
    -> std.atomic.load_i64
    -> std.sync.park
