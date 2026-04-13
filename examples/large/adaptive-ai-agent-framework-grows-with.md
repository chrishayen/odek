# Requirement: "an adaptive AI agent framework with tools, memory, and a planning loop that grows over time"

Agents carry a tool registry, a memory store, and a step loop that calls a pluggable model driver. Nothing in the library talks to a specific model provider.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a flat JSON object
      - returns error on malformed JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a flat JSON object
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

agent
  agent.new
    @ (name: string, system_prompt: string, model_fn: model_call_fn) -> agent_state
    + creates an agent bound to an injected model call function
    # construction
  agent.tool_register
    @ (state: agent_state, name: string, description: string, schema: string, run_fn: fn(string) -> result[string, string]) -> result[agent_state, string]
    + registers a tool the agent may invoke during a step
    - returns error when a tool with the same name is already registered
    # tools
  agent.tool_unregister
    @ (state: agent_state, name: string) -> result[agent_state, string]
    + removes a previously registered tool
    - returns error when no tool with that name exists
    # tools
  agent.memory_append
    @ (state: agent_state, role: string, content: string) -> agent_state
    + appends a message to the short-term conversation memory
    -> std.time.now_millis
    # memory
  agent.memory_snapshot
    @ (state: agent_state) -> list[memory_entry]
    + returns the current short-term memory as an ordered list
    # memory
  agent.long_term_store
    @ (state: agent_state, key: string, value: string) -> agent_state
    + writes a key/value pair to the long-term memory
    ? long-term memory persists across steps and is searchable by key
    # memory
  agent.long_term_recall
    @ (state: agent_state, key: string) -> optional[string]
    + returns a previously stored long-term value
    - returns none when the key is not present
    # memory
  agent.step
    @ (state: agent_state, user_input: string) -> result[step_result, string]
    + runs one reason-act cycle: builds a prompt, calls the model, parses tool calls, executes them, and appends results to memory
    - returns error when the model call fails
    - returns error when the model requests a tool that is not registered
    ? tool calls are detected by parsing a JSON "tool" field in the model response
    # step
    -> std.json.parse_object
    -> std.json.encode_object
    -> agent.memory_append
  agent.run_until
    @ (state: agent_state, user_input: string, stop_fn: fn(step_result) -> bool, max_steps: i32) -> result[step_result, string]
    + repeats step until stop_fn returns true or max_steps is reached
    - returns error when the max step budget is exhausted without stopping
    # loop
    -> agent.step
  agent.reflect
    @ (state: agent_state) -> result[agent_state, string]
    + asks the model to summarize recent memory and store the summary in long-term memory
    + intended to be called periodically so memory does not grow unbounded
    # reflection
    -> agent.long_term_store
    -> agent.memory_snapshot
  agent.learn_skill
    @ (state: agent_state, name: string, description: string, steps: list[string]) -> agent_state
    + stores a named procedure that the agent may later recall as a higher-level tool
    # adaptation
