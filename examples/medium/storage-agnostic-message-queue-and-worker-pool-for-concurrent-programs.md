# Requirement: "a storage-agnostic message queue and worker pool for concurrent programs"

The queue delegates persistence to a backend interface so callers can plug in memory, disk, or a database. The worker pool pulls jobs and dispatches them to handlers.

std: (all units exist)

jobqueue
  jobqueue.new_queue
    fn (backend: queue_backend) -> queue
    + wraps a backend in a queue handle
    # construction
  jobqueue.enqueue
    fn (q: queue, topic: string, payload: bytes) -> result[string, string]
    + returns the new job id
    - returns error when the backend rejects the write
    # enqueue
  jobqueue.dequeue
    fn (q: queue, topic: string) -> result[optional[job], string]
    + returns the next job for the topic, or none when empty
    - returns error when the backend fails
    # dequeue
  jobqueue.ack
    fn (q: queue, id: string) -> result[void, string]
    + marks the job complete and removes it from the backend
    - returns error when the id is unknown
    # ack
  jobqueue.nack
    fn (q: queue, id: string, requeue: bool) -> result[void, string]
    + returns the job to the backend or moves it to a dead-letter state
    - returns error when the id is unknown
    # nack
  jobqueue.new_pool
    fn (q: queue, topic: string, size: i32, handler: job_handler) -> pool
    + constructs a worker pool for the topic
    - returns an unstarted pool when size <= 0
    # pool_construction
  jobqueue.pool_tick
    fn (p: pool) -> result[i32, string]
    + pulls up to size jobs, runs each handler, and acks successes
    + returns the number of jobs processed in this tick
    - nacks any job whose handler returned an error
    # pool_work
  jobqueue.stop
    fn (p: pool) -> void
    + marks the pool stopped so subsequent ticks are no-ops
    # pool_lifecycle
