# Requirement: "a job submission client for cluster schedulers"

A uniform client for submitting and tracking batch jobs on a compute cluster. Resource specification, submission, and monitoring are the three concerns.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

cluster
  cluster.new_job_spec
    fn (name: string, command: string) -> job_spec
    + creates a minimal spec with defaults (1 cpu, default memory)
    # construction
  cluster.with_resources
    fn (spec: job_spec, cpus: i32, memory_mb: i64) -> job_spec
    + returns a spec with cpu count and memory requirement set
    - returns spec unchanged when cpus or memory are non-positive
    # configuration
  cluster.with_environment
    fn (spec: job_spec, env: map[string, string]) -> job_spec
    + returns a spec with environment variables applied
    # configuration
  cluster.submit
    fn (session: cluster_session, spec: job_spec) -> result[string, string]
    + submits the job and returns the cluster-assigned job id
    - returns error when the session is not connected
    # submission
    -> std.time.now_seconds
  cluster.status
    fn (session: cluster_session, job_id: string) -> result[job_state, string]
    + returns current state (queued, running, done, failed)
    - returns error when job_id is unknown
    # monitoring
  cluster.wait
    fn (session: cluster_session, job_id: string, timeout_s: i64) -> result[job_state, string]
    + polls until the job reaches a terminal state or timeout elapses
    - returns error when timeout is reached before completion
    # monitoring
    -> std.time.now_seconds
  cluster.cancel
    fn (session: cluster_session, job_id: string) -> result[void, string]
    + requests cancellation of a queued or running job
    - returns error when the job has already finished
    # control
