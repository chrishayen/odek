# Requirement: "a framework for building applications powered by language models"

The framework composes prompt templates, language model calls, tool invocations, memory, and chains. Model backends are pluggable by identifier.

std
  std.http
    std.http.post
      fn (url: string, headers: map[string, string], body: bytes) -> result[bytes, string]
      + sends a POST request and returns the body
      - returns error on non-2xx status
      # http
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on invalid input
      # serialization
  std.regex
    std.regex.replace_all
      fn (pattern: string, input: string, replacement: string) -> string
      + replaces every regex match
      # regex

llm_framework
  llm_framework.new_template
    fn (body: string, variable_names: list[string]) -> prompt_template
    + constructs a prompt template with named variables
    # prompting
  llm_framework.render_template
    fn (template: prompt_template, values: map[string, string]) -> result[string, string]
    + substitutes variable values into the template
    - returns error when a required variable is missing
    # prompting
    -> std.regex.replace_all
  llm_framework.register_model
    fn (name: string, endpoint: string, api_key: string) -> void
    + registers a language model backend
    # model_registration
  llm_framework.call_model
    fn (model_name: string, prompt: string) -> result[string, string]
    + sends the prompt to the named model and returns the completion text
    - returns error when the model is not registered
    - returns error on transport failure
    # model_invocation
    -> std.json.encode_object
    -> std.http.post
    -> std.json.parse_object
  llm_framework.register_tool
    fn (name: string, description: string, handler_id: string) -> void
    + registers a callable tool by name
    # tool_registration
  llm_framework.invoke_tool
    fn (name: string, arguments: map[string, string]) -> result[string, string]
    + invokes a registered tool and returns its result
    - returns error when the tool is not registered
    # tool_invocation
  llm_framework.new_memory
    fn () -> memory_state
    + creates an empty conversation memory
    # memory
  llm_framework.remember
    fn (state: memory_state, role: string, content: string) -> memory_state
    + appends a turn to the memory
    # memory
  llm_framework.memory_as_prompt
    fn (state: memory_state) -> string
    + renders the memory as a role-prefixed transcript
    # memory
  llm_framework.new_chain
    fn () -> chain_state
    + creates an empty chain
    # chain
  llm_framework.chain_prompt
    fn (state: chain_state, template: prompt_template, model_name: string) -> chain_state
    + appends a prompt-and-model step to the chain
    # chain
  llm_framework.chain_tool
    fn (state: chain_state, tool_name: string) -> chain_state
    + appends a tool-invocation step to the chain
    # chain
  llm_framework.run_chain
    fn (chain: chain_state, initial: map[string, string]) -> result[string, string]
    + executes each step in order, threading the previous output into the next
    - returns error on any step failure
    # chain
