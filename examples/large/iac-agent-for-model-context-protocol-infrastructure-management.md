# Requirement: "a tool that lets an AI assistant manage infrastructure-as-code environments over a model-context protocol"

Exposes plan/apply/destroy operations to an external model via a structured request/response protocol. The protocol codec and subprocess runner live in std; the project defines the tool methods and dispatch.

std
  std.process
    std.process.run
      @ (cmd: string, args: list[string], cwd: string) -> result[process_output, string]
      + captures stdout, stderr, and exit code
      - returns error when the binary is not found
      # process
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file's full contents as a string
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes contents atomically, creating parents if needed
      - returns error when the path is not writable
      # filesystem
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses any JSON document into a dynamic value
      - returns error on malformed JSON
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + serializes a dynamic value to a compact JSON string
      # serialization

iac_agent
  iac_agent.new
    @ (workspace_dir: string, binary: string) -> iac_agent_state
    + creates an agent rooted at a workspace with the given IaC binary
    # construction
  iac_agent.plan
    @ (state: iac_agent_state) -> result[string, string]
    + runs plan in the workspace and returns the human-readable diff
    - returns error when the workspace has no configuration file
    # plan
    -> std.process.run
  iac_agent.apply
    @ (state: iac_agent_state, auto_approve: bool) -> result[string, string]
    + runs apply and returns the output log
    - returns error when apply fails mid-run
    # apply
    -> std.process.run
  iac_agent.destroy
    @ (state: iac_agent_state) -> result[string, string]
    + runs destroy and returns the output log
    # destroy
    -> std.process.run
  iac_agent.read_config
    @ (state: iac_agent_state, relative_path: string) -> result[string, string]
    + returns the contents of a configuration file inside the workspace
    - returns error when the path escapes the workspace
    # config_read
    -> std.fs.read_all
  iac_agent.write_config
    @ (state: iac_agent_state, relative_path: string, contents: string) -> result[void, string]
    + writes a configuration file inside the workspace
    - returns error when the path escapes the workspace
    # config_write
    -> std.fs.write_all
  iac_agent.handle_request
    @ (state: iac_agent_state, request_json: string) -> string
    + decodes a protocol request, dispatches to the matching method, and encodes the response
    + unknown methods produce a structured error response
    # protocol_dispatch
    -> std.json.parse
    -> std.json.encode
