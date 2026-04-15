# Requirement: "a framework for IoT edge applications and integration flows"

A flow runtime where triggers fire actions composed of activities connected into a directed graph.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged tree
      - returns error on malformed JSON
      # serialization
    std.json.encode
      fn (value: json_value) -> string
      + serializes a tagged tree to JSON
      # serialization
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.v4
      fn () -> string
      + returns a random UUIDv4 string
      # uuid

flow_runtime
  flow_runtime.new
    fn () -> runtime_state
    + creates an empty runtime with no registered triggers or activities
    # construction
  flow_runtime.register_activity
    fn (rt: runtime_state, name: string, activity_id: string) -> result[runtime_state, string]
    + makes an activity available for use in flows
    - returns error when an activity with the same name is already registered
    # activity
  flow_runtime.register_trigger
    fn (rt: runtime_state, name: string, trigger_id: string) -> result[runtime_state, string]
    + makes a trigger type available for flows to bind to
    - returns error when the name is already registered
    # trigger
  flow_runtime.load_flow
    fn (rt: runtime_state, definition_json: string) -> result[tuple[string, runtime_state], string]
    + parses a flow definition (nodes, edges, start node) and assigns a flow id
    - returns error when JSON is malformed
    - returns error when the graph references unknown activities
    - returns error when the graph has a cycle
    # flow_loading
    -> std.json.parse
    -> std.uuid.v4
  flow_runtime.bind_trigger
    fn (rt: runtime_state, flow_id: string, trigger_name: string, config: map[string, string]) -> result[runtime_state, string]
    + attaches a configured trigger to a flow
    - returns error when flow_id or trigger_name is unknown
    # trigger
  flow_runtime.fire
    fn (rt: runtime_state, flow_id: string, input: map[string, string]) -> result[tuple[string, runtime_state], string]
    + creates a new execution instance for the flow with the given input
    + returns the execution id
    - returns error when flow_id is unknown
    # execution
    -> std.uuid.v4
    -> std.time.now_millis
  flow_runtime.step
    fn (rt: runtime_state, execution_id: string) -> result[tuple[execution_status, runtime_state], string]
    + runs the next ready node; returns running, waiting, completed, or failed
    - returns error when execution_id is unknown
    # execution
    -> std.time.now_millis
  flow_runtime.deliver_activity_result
    fn (rt: runtime_state, execution_id: string, node_id: string, output: map[string, string]) -> result[runtime_state, string]
    + records an activity's output and marks its successors as ready
    - returns error when node is not waiting for a result
    # execution
  flow_runtime.execution_status
    fn (rt: runtime_state, execution_id: string) -> optional[execution_status]
    + returns the current status of an execution
    # introspection
  flow_runtime.execution_output
    fn (rt: runtime_state, execution_id: string) -> optional[map[string, string]]
    + returns the final output of a completed execution
    # introspection
  flow_runtime.list_flows
    fn (rt: runtime_state) -> list[string]
    + returns the ids of every loaded flow
    # introspection
  flow_runtime.export_flow
    fn (rt: runtime_state, flow_id: string) -> result[string, string]
    + returns the flow definition serialized as JSON
    - returns error when flow_id is unknown
    # flow_export
    -> std.json.encode
