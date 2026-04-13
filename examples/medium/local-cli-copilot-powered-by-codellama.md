# Requirement: "a local code-assistant client that sends prompts to a language model endpoint"

Wraps prompt construction and request/response handling for a locally hosted completion model.

std
  std.http
    std.http.post_json
      @ (url: string, body: string) -> result[string, string]
      + posts a JSON body and returns the raw response body
      - returns error on network failure or non-2xx status
      # http
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as a JSON object
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

copilot
  copilot.new
    @ (endpoint: string, model: string) -> copilot_state
    + stores the endpoint URL and model identifier to use for subsequent requests
    # construction
  copilot.build_prompt
    @ (instruction: string, context: string) -> string
    + composes an instruction and surrounding code context into a single prompt string
    # prompt_construction
  copilot.complete
    @ (state: copilot_state, prompt: string) -> result[string, string]
    + sends the prompt to the configured endpoint and returns the generated text
    - returns error when the endpoint is unreachable or returns an error payload
    # inference
    -> std.json.encode_object
    -> std.http.post_json
    -> std.json.parse_object
  copilot.explain
    @ (state: copilot_state, code: string) -> result[string, string]
    + asks the model to produce a natural-language explanation of the given code
    # high_level_task
