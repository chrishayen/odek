# Requirement: "a dag workflow executor that runs steps defined in a simple declarative format"

Parses a workflow document, validates the graph, runs steps in dependency order, and tracks status.

std
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[yaml_value, string]
      + parses a declarative document into a generic value tree
      - returns error on malformed input
      # parsing
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

workflow
  workflow.parse
    @ (raw: string) -> result[workflow_def, string]
    + extracts name, steps, commands, and dependencies from a declarative document
    - returns error when required fields are missing
    # parsing
    -> std.yaml.parse
  workflow.validate
    @ (def: workflow_def) -> result[void, string]
    + checks that all dependencies reference known steps
    - returns error when a step references an unknown dependency
    - returns error when the graph contains a cycle
    - returns error when two steps share the same name
    # validation
  workflow.topological_order
    @ (def: workflow_def) -> result[list[string], string]
    + returns step names in dependency-first order
    - returns error when the graph contains a cycle
    # scheduling
  workflow.new_run
    @ (def: workflow_def) -> run_state
    + creates a run with all steps in pending status
    # execution
    -> std.time.now_millis
  workflow.ready_steps
    @ (run: run_state) -> list[string]
    + returns steps whose dependencies have all succeeded and are not yet started
    # execution
  workflow.mark_started
    @ (run: run_state, step: string) -> result[run_state, string]
    + transitions a step from pending to running
    - returns error when the step is unknown
    - returns error when the step is not pending
    # execution
    -> std.time.now_millis
  workflow.mark_succeeded
    @ (run: run_state, step: string) -> result[run_state, string]
    + transitions a step from running to succeeded
    - returns error when the step is not running
    # execution
    -> std.time.now_millis
  workflow.mark_failed
    @ (run: run_state, step: string, reason: string) -> result[run_state, string]
    + transitions a step from running to failed with a reason
    + marks all dependents as skipped
    - returns error when the step is not running
    # execution
    -> std.time.now_millis
  workflow.is_complete
    @ (run: run_state) -> bool
    + returns true when no step is pending or running
    # inspection
  workflow.status
    @ (run: run_state, step: string) -> result[step_status, string]
    + returns the current status of a step
    - returns error when the step is unknown
    # inspection
  workflow.summary
    @ (run: run_state) -> run_summary
    + returns counts of succeeded, failed, skipped, and pending steps
    # inspection
