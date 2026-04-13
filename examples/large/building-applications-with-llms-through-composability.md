# Requirement: "a library for building applications composed of language-model calls, tools, and retrieval"

Provides the abstractions for prompts, chains, tools, memory, and retrievers.

std
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on invalid JSON or non-object root
      # serialization
  std.strings
    std.strings.replace_all
      @ (s: string, needle: string, replacement: string) -> string
      + replaces every occurrence
      # strings

llmkit
  llmkit.prompt_template
    @ (template: string) -> prompt_template
    + creates a template containing {name} placeholders
    # prompting
  llmkit.render_prompt
    @ (tmpl: prompt_template, vars: map[string, string]) -> result[string, string]
    + substitutes variables into the template
    - returns error when a placeholder has no binding
    # prompting
    -> std.strings.replace_all
  llmkit.llm_call
    @ (endpoint: callback, prompt: string) -> result[string, string]
    + invokes an endpoint function with the rendered prompt and returns the completion
    - returns error when the endpoint reports failure
    # model_invocation
  llmkit.chain
    @ (steps: list[chain_step]) -> chain_state
    + composes steps whose outputs feed the next step's inputs
    # composition
  llmkit.run_chain
    @ (chain: chain_state, inputs: map[string, string]) -> result[map[string, string], string]
    + executes steps in order, threading outputs forward
    - returns error when any step fails
    # composition
  llmkit.register_tool
    @ (name: string, description: string, handler: callback) -> tool_def
    + declares a tool the model can invoke with a JSON argument
    # tools
  llmkit.parse_tool_call
    @ (text: string) -> result[tuple[string, map[string, string]], string]
    + extracts tool name and argument map from model output
    - returns error on malformed tool call
    # tools
    -> std.json.parse_object
  llmkit.run_tool
    @ (tools: list[tool_def], name: string, args: map[string, string]) -> result[string, string]
    + finds the named tool and runs its handler with the arguments
    - returns error when the tool is not registered
    # tools
    -> std.json.encode_object
  llmkit.memory_new
    @ () -> memory_state
    + creates empty conversation memory
    # memory
  llmkit.memory_append
    @ (mem: memory_state, role: string, content: string) -> memory_state
    + appends a role-tagged turn
    # memory
  llmkit.memory_format
    @ (mem: memory_state) -> string
    + serializes memory as a chat transcript suitable for prompt injection
    # memory
  llmkit.retriever_add
    @ (store: retriever_state, doc_id: string, embedding: list[f32], text: string) -> retriever_state
    + indexes a document by its embedding
    # retrieval
  llmkit.retriever_search
    @ (store: retriever_state, query_embedding: list[f32], k: i32) -> list[retrieval_hit]
    + returns the k documents with highest cosine similarity
    # retrieval
