# Requirement: "a software transactional memory library"

Optimistic transactions over versioned cells: read your own writes, commit atomically, retry on conflict.

std
  std.sync
    std.sync.new_mutex
      @ () -> mutex
      + creates an unlocked mutex
      # concurrency
    std.sync.lock
      @ (m: mutex) -> void
      + acquires the mutex, blocking if held
      # concurrency
    std.sync.unlock
      @ (m: mutex) -> void
      + releases a held mutex
      # concurrency

stm
  stm.new_cell
    @ (initial: i64) -> cell
    + allocates a versioned cell with the given initial value
    # construction
    -> std.sync.new_mutex
  stm.begin
    @ () -> tx
    + starts a fresh transaction with empty read and write sets
    # transaction_lifecycle
  stm.read
    @ (t: tx, c: cell) -> i64
    + returns the current value, recording the cell and its version in the read set
    ? subsequent reads within the same tx see the tx's own writes
    # transaction_read
  stm.write
    @ (t: tx, c: cell, value: i64) -> void
    + buffers a pending write; visible only to this transaction until commit
    # transaction_write
  stm.commit
    @ (t: tx) -> bool
    + validates the read set against current versions and applies writes atomically
    + returns true on success and false when any read-set version has advanced
    - returns false when a concurrent tx already committed a conflicting write
    # commit
    -> std.sync.lock
    -> std.sync.unlock
  stm.atomically
    @ (body: tx_fn) -> void
    + runs body inside a transaction, retrying with fresh state until commit succeeds
    ? the body must be idempotent; it may run multiple times
    # retry_loop
