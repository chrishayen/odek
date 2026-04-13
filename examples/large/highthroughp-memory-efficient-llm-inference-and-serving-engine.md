# Requirement: "a high-throughput and memory-efficient inference and serving engine for language models"

The core ideas are a paged key-value cache, a request scheduler that batches compatible requests, and a serving loop that advances batches one step at a time.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.sync
    std.sync.new_mutex
      @ () -> mutex_handle
      + returns a new unlocked mutex handle
      # concurrency
    std.sync.lock
      @ (m: mutex_handle) -> void
      + blocks until the mutex is acquired
      # concurrency
    std.sync.unlock
      @ (m: mutex_handle) -> void
      + releases a mutex held by the current caller
      - panics if the mutex is not held
      # concurrency
  std.collections
    std.collections.new_priority_queue
      @ () -> pqueue_state
      + returns an empty min-heap priority queue
      # collections
    std.collections.pqueue_push
      @ (q: pqueue_state, priority: i64, item_id: i64) -> pqueue_state
      + inserts an item with a priority key
      # collections
    std.collections.pqueue_pop
      @ (q: pqueue_state) -> tuple[optional[i64], pqueue_state]
      + returns the item with the smallest priority and the new queue
      - returns none when the queue is empty
      # collections

inference_engine
  inference_engine.new_kv_cache
    @ (num_blocks: i32, block_size: i32) -> kv_cache_state
    + allocates a paged key-value cache with the given block count and block size
    ? blocks are opaque slot ids; actual tensor storage is host-provided
    # memory_management
  inference_engine.allocate_blocks
    @ (cache: kv_cache_state, num_tokens: i32) -> result[tuple[list[i32], kv_cache_state], string]
    + returns the block ids reserved to hold num_tokens and the updated cache
    - returns error when not enough free blocks remain
    # memory_management
  inference_engine.free_blocks
    @ (cache: kv_cache_state, block_ids: list[i32]) -> kv_cache_state
    + returns the listed blocks to the free pool
    # memory_management
  inference_engine.new_scheduler
    @ (max_batch_tokens: i32) -> scheduler_state
    + creates a scheduler bounded by a token budget per batch
    # scheduling
  inference_engine.add_request
    @ (sched: scheduler_state, request_id: i64, prompt_tokens: list[i32], max_new: i32) -> scheduler_state
    + queues a new generation request with its prompt and max new tokens
    # scheduling
    -> std.time.now_millis
  inference_engine.build_batch
    @ (sched: scheduler_state, cache: kv_cache_state) -> tuple[list[i64], scheduler_state, kv_cache_state]
    + selects a set of requests that fit the token budget and have blocks available
    + returns the batched request ids and updated scheduler/cache state
    # scheduling
    -> std.collections.pqueue_pop
  inference_engine.step_batch
    @ (sched: scheduler_state, request_ids: list[i64], sampled_tokens: list[i32]) -> tuple[list[i64], scheduler_state]
    + appends one sampled token per request and returns the ids that finished this step
    ? the actual forward pass is caller-provided; this unit manages bookkeeping
    # generation
  inference_engine.get_output
    @ (sched: scheduler_state, request_id: i64) -> optional[list[i32]]
    + returns all tokens generated so far for a request
    - returns none for an unknown request id
    # output
  inference_engine.cancel_request
    @ (sched: scheduler_state, request_id: i64) -> scheduler_state
    + removes a request and frees its cache blocks
    # scheduling
    -> std.collections.new_priority_queue
