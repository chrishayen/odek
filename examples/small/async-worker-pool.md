# Requirement: "a simple async worker pool"

A fixed-size pool that accepts job closures on a queue and runs them concurrently. The pool is modeled as opaque state; the host runtime provides the actual goroutines via a thin spawn primitive.

std
  std.concurrency
    std.concurrency.spawn
      fn (work: closure[void]) -> void
      + runs the closure on a host-managed thread of execution
      # concurrency

worker_pool
  worker_pool.new
    fn (worker_count: i32, queue_capacity: i32) -> pool_state
    + creates a pool with the given worker count and bounded queue
    ? worker goroutines are started lazily on the first submit
    # construction
  worker_pool.submit
    fn (state: pool_state, job: closure[void]) -> result[void, string]
    + enqueues the job for execution by an available worker
    - returns error when the queue is full and the pool is non-blocking
    - returns error when the pool has been shut down
    # submission
    -> std.concurrency.spawn
  worker_pool.shutdown
    fn (state: pool_state) -> void
    + signals workers to drain the queue and exit; blocks until all finish
    # lifecycle
