# Requirement: "a micro-batching library"

Groups individual inputs into fixed-size or time-bounded batches and flushes them to a caller-provided processor.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

micro_batch
  micro_batch.new
    @ (max_size: i32, max_wait_ms: i64) -> batcher_state
    + creates a batcher with size and latency thresholds
    ? max_size must be positive
    # construction
  micro_batch.submit
    @ (state: batcher_state, item: bytes) -> tuple[batcher_state, optional[list[bytes]]]
    + returns (new_state, some_batch) when the max_size threshold is reached, otherwise (new_state, none)
    # submission
    -> std.time.now_millis
  micro_batch.tick
    @ (state: batcher_state) -> tuple[batcher_state, optional[list[bytes]]]
    + returns a batch when the pending items have been waiting longer than max_wait_ms
    - returns (state, none) when the batch is still within the wait window
    # flushing
    -> std.time.now_millis
  micro_batch.flush
    @ (state: batcher_state) -> tuple[batcher_state, list[bytes]]
    + returns all pending items as a final batch and empties the buffer
    # shutdown
