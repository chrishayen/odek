# Requirement: "a lightweight asynchronous queue job manager"

A minimal in-memory queue with background execution. No broker, no persistence.

std: (all units exist)

asyncjob
  asyncjob.manager_new
    @ (concurrency: i32) -> manager_state
    + creates a manager with the given worker count
    ? concurrency<=0 is clamped to 1
    # construction
  asyncjob.submit
    @ (manager: manager_state, job: job_fn) -> i64
    + enqueues a job and returns its monotonically increasing id
    # submission
  asyncjob.start
    @ (manager: manager_state) -> manager_state
    + spawns worker tasks that drain the queue
    # lifecycle
  asyncjob.stop
    @ (manager: manager_state) -> void
    + signals workers to finish their current job and exit
    # lifecycle
  asyncjob.status
    @ (manager: manager_state, job_id: i64) -> optional[job_status]
    + returns pending, running, done, or failed for a job id
    - returns none for an unknown id
    # inspection
