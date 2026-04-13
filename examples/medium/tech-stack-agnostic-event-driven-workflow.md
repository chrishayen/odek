# Requirement: "an event-driven workflow framework"

Workflows are declared as steps keyed by the event that triggers them. Dispatching an event advances any workflow instance waiting on it.

std: (all units exist)

workflow
  workflow.new
    @ () -> workflow_engine
    + creates an empty engine with no registered steps
    # construction
  workflow.register_step
    @ (engine: workflow_engine, event: string, step_id: string, handler: step_fn) -> workflow_engine
    + registers a handler to run when a workflow instance is waiting on event at step_id
    ? step_fn takes the current state and event payload and returns (next_step_id, new_state)
    # registration
  workflow.start_instance
    @ (engine: workflow_engine, initial_step: string, initial_state: map[string,string]) -> tuple[string, workflow_engine]
    + creates an instance in initial_step and returns its id
    # instance_lifecycle
  workflow.dispatch
    @ (engine: workflow_engine, instance_id: string, event: string, payload: map[string,string]) -> result[workflow_engine, string]
    + runs the handler registered for (current_step, event) and advances the instance
    - returns error when no step is registered for (current_step, event)
    - returns error when instance_id is unknown
    # dispatch
  workflow.instance_state
    @ (engine: workflow_engine, instance_id: string) -> result[map[string,string], string]
    + returns the current state map for an instance
    - returns error when instance_id is unknown
    # inspection
