# Requirement: "a local ai agent that automates engineering tasks"

Agent loop: call a pluggable language model, dispatch tool invocations, iterate until done. Model and tools are injected — no specific provider.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the file does not exist
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes, creating the file if needed
      # io
  std.proc
    std.proc.run
      @ (cmd: string, args: list[string]) -> result[proc_result, string]
      + runs the command and returns stdout, stderr, exit code
      - returns error when the binary is not found
      # process
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on malformed input
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

agent
  agent.new_session
    @ (system_prompt: string) -> session_state
    + returns a session seeded with the system prompt and empty history
    # session
  agent.append_user_message
    @ (s: session_state, text: string) -> session_state
    + appends a user message to the session
    # session
  agent.append_assistant_message
    @ (s: session_state, text: string) -> session_state
    + appends an assistant message to the session
    # session
  agent.register_tool
    @ (s: session_state, name: string, description: string) -> session_state
    + adds a tool descriptor the model may call
    - rejects a name already registered
    # tools
  agent.parse_tool_call
    @ (assistant_text: string) -> optional[tool_invocation]
    + returns the parsed tool invocation when the assistant text contains one
    + returns none when the assistant did not request a tool
    # tools
    -> std.json.parse_object
  agent.dispatch_tool
    @ (s: session_state, invocation: tool_invocation) -> result[string, string]
    + runs the registered tool and returns its output
    - returns error when the tool name is not registered
    # tools
    -> std.proc.run
    -> std.fs.read_all
    -> std.fs.write_all
  agent.tool_result_message
    @ (name: string, output: string) -> string
    + formats a tool result for the next turn
    -> std.json.encode_object
    # tools
  agent.step
    @ (s: session_state, model_call: model_fn) -> result[session_state, string]
    + calls the model, appends its reply, executes any tool call, and appends the tool result
    ? a single turn of the loop; caller decides when to stop
    # orchestration
  agent.is_done
    @ (s: session_state) -> bool
    + returns true when the last assistant message contains no tool call
    # orchestration
