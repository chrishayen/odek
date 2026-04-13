# Requirement: "a configurable, expressive parallel computation library"

A pipeline builder that describes data-parallel computations over input sequences. Map, filter, and reduce stages execute against a configurable worker pool.

std
  std.concurrency
    std.concurrency.worker_pool_new
      @ (num_workers: i32) -> pool_state
      + returns a pool with the given number of workers
      ? num_workers <= 0 falls back to 1
      # concurrency
    std.concurrency.submit_batch
      @ (pool: pool_state, tasks: list[bytes]) -> list[bytes]
      + runs tasks across workers and returns results in input order
      # concurrency

parallel
  parallel.from_list
    @ (items: list[bytes]) -> pipeline_state
    + wraps a list of opaque items as a pipeline source
    # construction
  parallel.map
    @ (pipeline: pipeline_state, fn_id: string) -> pipeline_state
    + appends a map stage that applies the named function to each item
    # stages
  parallel.filter
    @ (pipeline: pipeline_state, predicate_id: string) -> pipeline_state
    + appends a filter stage that keeps items where the predicate returns true
    # stages
  parallel.reduce
    @ (pipeline: pipeline_state, reducer_id: string, initial: bytes) -> pipeline_state
    + appends a reduce stage producing a single output
    # stages
  parallel.with_workers
    @ (pipeline: pipeline_state, num_workers: i32) -> pipeline_state
    + configures the number of worker threads for execution
    # configuration
  parallel.with_chunk_size
    @ (pipeline: pipeline_state, chunk_size: i32) -> pipeline_state
    + sets how many items each worker processes per task
    ? larger chunks reduce scheduling overhead but increase latency variance
    # configuration
  parallel.run
    @ (pipeline: pipeline_state) -> result[list[bytes], string]
    + executes the pipeline and returns the final items
    - returns error when an unknown fn_id is referenced
    # execution
    -> std.concurrency.worker_pool_new
    -> std.concurrency.submit_batch
