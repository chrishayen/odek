# Requirement: "async utilities for managing closures and callbacks"

Small combinator library: run callbacks in parallel, in sequence, or with fallbacks.

std: (all units exist)

gollback
  gollback.all
    fn (tasks: list[task_fn]) -> list[result[any, string]]
    + runs all tasks concurrently and returns every outcome in order
    + preserves input ordering regardless of completion order
    # parallel
  gollback.race
    fn (tasks: list[task_fn]) -> result[any, string]
    + returns the first successful result and cancels the rest
    - returns the last error when all tasks fail
    # racing
  gollback.retry
    fn (task: task_fn, max_attempts: i32) -> result[any, string]
    + invokes task until it succeeds or max_attempts is reached
    - returns the final error when every attempt fails
    # retry
