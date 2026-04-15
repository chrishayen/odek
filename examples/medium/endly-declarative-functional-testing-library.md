# Requirement: "a declarative end-to-end functional testing library"

A workflow document describes ordered steps, each dispatched to a registered action; the runner collects per-step results and decides pass/fail.

std: (all units exist)

endly
  endly.parse_workflow
    fn (source: string) -> result[workflow, string]
    + parses a document with ordered steps, each with name and action
    - returns error when a step has no action
    # parsing
  endly.register_action
    fn (name: string, handler: action_fn) -> void
    + registers an action handler under a name
    # registry
  endly.new_context
    fn (variables: map[string, string]) -> run_context
    + returns a context seeded with the given variables
    # construction
  endly.expand_template
    fn (ctx: run_context, template: string) -> string
    + substitutes ${name} references from context variables
    + leaves unknown references unchanged
    # templating
  endly.run_step
    fn (ctx: run_context, step: step) -> tuple[step_result, run_context]
    + dispatches to the step's action and returns its result plus the updated context
    - returns a failed result when the action name is unregistered
    # execution
    -> endly.expand_template
  endly.run_workflow
    fn (ctx: run_context, wf: workflow) -> workflow_result
    + runs steps in order and aggregates results
    + stops at the first failure when halt_on_error is set
    # orchestration
    -> endly.run_step
  endly.assert_result
    fn (result: workflow_result) -> result[void, string]
    - returns error listing failed step names when any step failed
    # assertion
