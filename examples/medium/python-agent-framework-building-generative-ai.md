# Requirement: "an agent framework for building generative AI applications with structured schemas"

Wires a model call loop around a schema validator and tool registry. The model call itself is a std seam so tests can substitute a deterministic responder.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.http
    std.http.post_json
      @ (url: string, headers: map[string, string], body: string) -> result[string, string]
      + posts a JSON body and returns the response body as a string
      - returns error on non-2xx status
      # network

agent
  agent.new
    @ (model_endpoint: string, api_key: string) -> agent_state
    + creates an agent configured with an LLM endpoint and credentials
    # construction
  agent.define_schema
    @ (state: agent_state, name: string, fields: map[string, string]) -> agent_state
    + registers a named output schema mapping field_name to type_name
    ? schemas coerce the model's JSON output before it reaches the caller
    # schema_registry
  agent.register_tool
    @ (state: agent_state, name: string, description: string, input_fields: map[string, string]) -> agent_state
    + adds a tool the agent can call during a run
    # tool_registry
  agent.run
    @ (state: agent_state, prompt: string, schema_name: string) -> result[map[string, string], string]
    + sends the prompt to the model, parses the response, validates against the named schema, and returns the fields
    - returns error when the schema_name is unknown
    - returns error when the model response is not valid JSON matching the schema
    # agent_loop
    -> std.http.post_json
    -> std.json.parse_object
    -> std.json.encode_object
  agent.validate_against_schema
    @ (state: agent_state, schema_name: string, obj: map[string, string]) -> result[void, string]
    + returns ok when every required field is present with a matching type
    - returns error listing missing or mistyped fields
    # validation
