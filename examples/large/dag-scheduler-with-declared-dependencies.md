# Requirement: "a DAG scheduler that runs tasks with declared dependencies"

Tasks form a directed acyclic graph; the scheduler executes them in topological order, respecting per-task state and reporting progress.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns the current time in milliseconds since the epoch
      # time
  std.concurrency
    std.concurrency.run_parallel
      fn (fns: list[task_fn], max_parallel: i32) -> list[result[string, string]]
      + invokes fns with a bounded concurrency level and returns each result in order
      # concurrency

dag_scheduler
  dag_scheduler.new
    fn () -> dag_state
    + creates an empty DAG with no tasks
    # construction
  dag_scheduler.add_task
    fn (dag: dag_state, id: string, fn: task_fn) -> result[dag_state, string]
    + registers a task with the given id and function
    - returns error when id already exists
    # construction
  dag_scheduler.add_dependency
    fn (dag: dag_state, from_id: string, to_id: string) -> result[dag_state, string]
    + records that to_id must run after from_id
    - returns error when either id is unknown
    - returns error when the edge would create a cycle
    # construction
  dag_scheduler.topological_order
    fn (dag: dag_state) -> result[list[string], string]
    + returns task ids in an order compatible with all dependencies
    - returns error when the graph contains a cycle
    # ordering
  dag_scheduler.ready_tasks
    fn (dag: dag_state) -> list[string]
    + returns ids whose dependencies are all marked completed
    # scheduling
  dag_scheduler.run
    fn (dag: dag_state, max_parallel: i32) -> run_report
    + executes tasks wave-by-wave, running each wave's ready tasks in parallel
    + records per-task status, start time, and duration in the report
    - marks a task as skipped when any dependency failed
    # execution
    -> std.concurrency.run_parallel
    -> std.time.now_millis
  dag_scheduler.task_status
    fn (report: run_report, id: string) -> optional[task_status]
    + returns the recorded status for a task id or none
    # reporting
  dag_scheduler.summary
    fn (report: run_report) -> run_summary
    + aggregates counts of completed, failed, and skipped tasks with total wall time
    # reporting
  dag_scheduler.render_text
    fn (report: run_report) -> string
    + returns a human-readable progress table keyed by task id
    # reporting
  dag_scheduler.reset
    fn (dag: dag_state) -> dag_state
    + clears all task statuses while preserving the graph structure
    # lifecycle
