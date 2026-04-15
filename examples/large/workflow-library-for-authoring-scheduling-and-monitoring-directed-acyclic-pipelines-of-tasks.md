# Requirement: "a workflow library for authoring, scheduling, and monitoring directed acyclic pipelines of tasks"

Pure library: define DAGs, compute runnable tasks, record run state. No execution, no UI.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
    std.time.parse_cron
      fn (expr: string) -> result[cron_spec, string]
      + parses a five-field cron expression
      - returns error on invalid cron syntax
      # time
    std.time.next_cron_fire
      fn (spec: cron_spec, after: i64) -> i64
      + returns the next fire time at or after the given unix second
      # time

workflow
  workflow.new_dag
    fn (name: string) -> dag_state
    + returns an empty DAG with the given name
    # construction
  workflow.add_task
    fn (dag: dag_state, task_id: string, tag: string) -> result[dag_state, string]
    + adds a task node identified by task_id
    - returns error when task_id already exists
    # authoring
  workflow.add_dependency
    fn (dag: dag_state, upstream: string, downstream: string) -> result[dag_state, string]
    + adds an edge from upstream to downstream
    - returns error when either task is unknown
    - returns error when the edge would create a cycle
    # authoring
  workflow.topological_order
    fn (dag: dag_state) -> result[list[string], string]
    + returns tasks in a valid execution order
    - returns error when the graph contains a cycle
    # analysis
  workflow.set_schedule
    fn (dag: dag_state, cron_expr: string) -> result[dag_state, string]
    + stores a cron schedule for the DAG
    - returns error on invalid cron syntax
    # scheduling
    -> std.time.parse_cron
  workflow.next_run_at
    fn (dag: dag_state, after: i64) -> optional[i64]
    + returns the next scheduled run time at or after the given unix second
    # scheduling
    -> std.time.next_cron_fire
  workflow.start_run
    fn (dag: dag_state, run_id: string) -> result[run_state, string]
    + returns a fresh run with all tasks queued
    - returns error when the DAG has a cycle
    # execution_state
    -> workflow.topological_order
    -> std.time.now_seconds
  workflow.mark_task
    fn (run: run_state, task_id: string, status: string) -> result[run_state, string]
    + updates the task status (queued, running, success, failed, skipped)
    - returns error when task_id is unknown
    - returns error when status is not one of the allowed values
    # execution_state
    -> std.time.now_seconds
  workflow.ready_tasks
    fn (run: run_state) -> list[string]
    + returns tasks whose upstream dependencies have all succeeded
    # execution_state
  workflow.run_status
    fn (run: run_state) -> string
    + returns one of queued, running, success, failed, partial
    # monitoring
  workflow.task_history
    fn (run: run_state, task_id: string) -> list[task_event]
    + returns the ordered status events for a task
    # monitoring
