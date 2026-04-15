# Requirement: "a framework for building agentic AI applications"

Defines agents with tools, a message conversation, and a run loop that dispatches tool calls to registered handlers.

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
  std.uuid
    std.uuid.v4
      fn () -> string
      + returns a random UUIDv4 string
      # uuid

agents
  agents.new_agent
    fn (name: string, system_prompt: string) -> agent_state
    + creates an agent with no tools and an empty conversation
    # construction
  agents.register_tool
    fn (agent: agent_state, name: string, description: string, parameters_schema: string) -> result[agent_state, string]
    + adds a tool definition to the agent
    - returns error when a tool with the same name is already registered
    - returns error when parameters_schema is not valid JSON schema
    # tool_registration
    -> std.json.parse
  agents.add_user_message
    fn (agent: agent_state, content: string) -> agent_state
    + appends a user message to the conversation
    # conversation
  agents.add_assistant_message
    fn (agent: agent_state, content: string, tool_calls: list[tool_call]) -> agent_state
    + appends an assistant message with optional tool calls
    # conversation
  agents.add_tool_result
    fn (agent: agent_state, tool_call_id: string, content: string) -> result[agent_state, string]
    - returns error when no matching tool call id is pending
    # conversation
  agents.pending_tool_calls
    fn (agent: agent_state) -> list[tool_call]
    + returns tool calls in the latest assistant message that have no result
    # introspection
  agents.dispatch
    fn (agent: agent_state, handlers: map[string, tool_handler]) -> result[agent_state, string]
    + invokes each registered handler for its matching pending call and records the result
    - returns error when a pending call has no handler
    # dispatch
    -> std.json.encode
    -> std.uuid.v4
  agents.build_request
    fn (agent: agent_state) -> string
    + returns a JSON request body for a model provider including system, tools, and messages
    # serialization
    -> std.json.encode
  agents.apply_model_response
    fn (agent: agent_state, response_json: string) -> result[agent_state, string]
    + parses a model response and appends the resulting assistant message
    - returns error on malformed response
    # deserialization
    -> std.json.parse
  agents.is_finished
    fn (agent: agent_state) -> bool
    + returns true when the latest assistant message has content and no pending tool calls
    # introspection
  agents.reset_conversation
    fn (agent: agent_state) -> agent_state
    + clears messages while preserving tools and system prompt
    # conversation
