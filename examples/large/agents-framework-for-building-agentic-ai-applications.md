# Requirement: "a framework for building agentic AI applications"

Defines agents with tools, a message conversation, and a run loop that dispatches tool calls to registered handlers.

std
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged tree
      - returns error on malformed JSON
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + serializes a tagged tree to JSON
      # serialization
  std.uuid
    std.uuid.v4
      @ () -> string
      + returns a random UUIDv4 string
      # uuid

agents
  agents.new_agent
    @ (name: string, system_prompt: string) -> agent_state
    + creates an agent with no tools and an empty conversation
    # construction
  agents.register_tool
    @ (agent: agent_state, name: string, description: string, parameters_schema: string) -> result[agent_state, string]
    + adds a tool definition to the agent
    - returns error when a tool with the same name is already registered
    - returns error when parameters_schema is not valid JSON schema
    # tool_registration
    -> std.json.parse
  agents.add_user_message
    @ (agent: agent_state, content: string) -> agent_state
    + appends a user message to the conversation
    # conversation
  agents.add_assistant_message
    @ (agent: agent_state, content: string, tool_calls: list[tool_call]) -> agent_state
    + appends an assistant message with optional tool calls
    # conversation
  agents.add_tool_result
    @ (agent: agent_state, tool_call_id: string, content: string) -> result[agent_state, string]
    - returns error when no matching tool call id is pending
    # conversation
  agents.pending_tool_calls
    @ (agent: agent_state) -> list[tool_call]
    + returns tool calls in the latest assistant message that have no result
    # introspection
  agents.dispatch
    @ (agent: agent_state, handlers: map[string, tool_handler]) -> result[agent_state, string]
    + invokes each registered handler for its matching pending call and records the result
    - returns error when a pending call has no handler
    # dispatch
    -> std.json.encode
    -> std.uuid.v4
  agents.build_request
    @ (agent: agent_state) -> string
    + returns a JSON request body for a model provider including system, tools, and messages
    # serialization
    -> std.json.encode
  agents.apply_model_response
    @ (agent: agent_state, response_json: string) -> result[agent_state, string]
    + parses a model response and appends the resulting assistant message
    - returns error on malformed response
    # deserialization
    -> std.json.parse
  agents.is_finished
    @ (agent: agent_state) -> bool
    + returns true when the latest assistant message has content and no pending tool calls
    # introspection
  agents.reset_conversation
    @ (agent: agent_state) -> agent_state
    + clears messages while preserving tools and system prompt
    # conversation
