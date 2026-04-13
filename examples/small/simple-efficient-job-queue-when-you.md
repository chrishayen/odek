# Requirement: "a simple in-process job queue"

An in-memory FIFO queue of named jobs with payloads, supporting enqueue, dequeue, and retry on failure.

std: (all units exist)

queue
  queue.new
    @ (max_retries: i32) -> queue_state
    + returns an empty queue that will retry each job up to max_retries times
    # construction
  queue.enqueue
    @ (state: queue_state, name: string, payload: bytes) -> queue_state
    + appends the job to the tail of the queue
    # producer
  queue.dequeue
    @ (state: queue_state) -> tuple[optional[job], queue_state]
    + returns the next job and removes it from the head
    - returns none when the queue is empty
    # consumer
  queue.mark_failed
    @ (state: queue_state, job: job) -> queue_state
    + re-enqueues the job at the tail with its attempt count incremented
    + drops the job permanently when attempts exceed max_retries
    # retry
  queue.pending_count
    @ (state: queue_state) -> i32
    + returns the number of jobs currently waiting
    # introspection
