# Requirement: "a library for building modular agents powered by large language models"

An agent is a named system prompt plus a set of tools; calling it streams messages to a pluggable model provider, optionally dispatching tool calls and feeding results back.

std
  std.http
    std.http.post_json
      @ (url: string, headers: map[string, string], body: string) -> result[string, string]
      + posts a JSON body and returns the response body on 2xx
      - returns error on network failure or non-2xx status
      # http
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged value
      - returns error on malformed input
      # serialization
    std.json.encode_value
      @ (v: json_value) -> string
      + encodes a tagged value as JSON
      # serialization
    std.json.get_string
      @ (v: json_value, path: list[string]) -> optional[string]
      + reads a string at a nested key path
      # serialization

agent
  agent.new_tool_registry
    @ () -> tool_registry
    + creates an empty tool registry
    # construction
  agent.register_tool
    @ (registry: tool_registry, name: string, schema: string, handler: tool_handler) -> result[tool_registry, string]
    + adds a tool with a JSON schema and a handler that runs it
    - returns error when the schema is not valid JSON
    - returns error when a tool with that name already exists
    # registration
    -> std.json.parse_value
  agent.new
    @ (system_prompt: string, tools: tool_registry, provider: model_provider) -> agent_state
    + creates an agent bound to a system prompt, tool set, and provider
    # construction
  agent.run
    @ (state: agent_state, user_message: string) -> result[list[agent_message], string]
    + returns the full exchange: user, assistant, tool calls, tool results
    + loops until the model returns a final assistant message with no tool calls
    - returns error when the provider returns a protocol error
    - returns error when a requested tool is not registered
    - returns error when the loop exceeds a safety limit of 25 steps
    # orchestration
  agent.build_request_body
    @ (state: agent_state, history: list[agent_message]) -> string
    + serializes the conversation and tool definitions for the provider
    # protocol
    -> std.json.encode_value
  agent.parse_response
    @ (raw: string) -> result[agent_step, string]
    + returns either a final assistant text or a list of tool calls
    - returns error on unexpected response shape
    # protocol
    -> std.json.parse_value
    -> std.json.get_string
  agent.dispatch_tool_call
    @ (tools: tool_registry, name: string, args: string) -> result[string, string]
    + invokes the named tool handler with the JSON arguments
    - returns error when the tool is unknown
    - returns error when the handler fails
    # tools
  agent.call_provider
    @ (provider: model_provider, body: string) -> result[string, string]
    + posts the request body to the configured endpoint with auth headers
    - returns error on transport or authorization failure
    # transport
    -> std.http.post_json
