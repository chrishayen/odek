# Requirement: "a workflow automation engine that models and executes business processes declared as YAML"

Loads a workflow definition, instantiates runs, executes steps with templated inputs, tracks status, and supports pausing for external approvals.

std
  std.yaml
    std.yaml.parse
      @ (raw: string) -> result[yaml_value, string]
      + parses YAML text into a generic value
      - returns error on malformed input
      # serialization
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + encodes a JSON value as text
      # serialization
    std.json.decode
      @ (raw: string) -> result[json_value, string]
      + parses text as a JSON value
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.ids
    std.ids.new_uuid
      @ () -> string
      + returns a fresh UUIDv4 string
      # identifiers

workflow
  workflow.parse_template
    @ (raw: string) -> result[template, string]
    + parses a YAML workflow definition into a typed template
    - returns error when a step references an unknown step as input
    # parsing
    -> std.yaml.parse
  workflow.validate_template
    @ (tpl: template) -> result[void, string]
    + returns ok when the template has no cycles and every input is resolvable
    - returns error on cyclic step dependencies
    # validation
  workflow.new_run
    @ (tpl: template, inputs: map[string, json_value]) -> run
    + creates a new run with a fresh id, the template, and the supplied inputs
    # execution
    -> std.ids.new_uuid
    -> std.time.now_millis
  workflow.register_action
    @ (reg: action_registry, name: string, action: action_fn) -> action_registry
    + adds a named action that steps can reference
    # registration
  workflow.resolve_inputs
    @ (step: step, ctx: run_context) -> map[string, json_value]
    + returns the step's input values after resolving references to prior step outputs
    # templating
  workflow.execute_step
    @ (reg: action_registry, step: step, inputs: map[string, json_value]) -> step_result
    + invokes the step's action and returns its outputs and final status
    - returns a failing result when the action returns an error
    # execution
  workflow.tick
    @ (r: run, reg: action_registry) -> run
    + runs all steps whose dependencies are satisfied and that are not blocked on approval
    + updates step statuses in place
    # execution
    -> std.time.now_millis
  workflow.request_approval
    @ (r: run, step_name: string, approver: string) -> run
    + marks a step as waiting for an external approval token
    # approvals
    -> std.ids.new_uuid
  workflow.resolve_approval
    @ (r: run, token: string, decision: string) -> result[run, string]
    + advances the step past approval when decision is "approve"
    - returns error when the token is unknown
    # approvals
  workflow.run_status
    @ (r: run) -> string
    + returns one of "running", "waiting", "succeeded", "failed"
    # status
  workflow.serialize_run
    @ (r: run) -> string
    + returns a JSON representation of the run for persistence
    # persistence
    -> std.json.encode
  workflow.load_run
    @ (raw: string) -> result[run, string]
    + restores a run from its serialized form
    - returns error on malformed JSON
    # persistence
    -> std.json.decode
