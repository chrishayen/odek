# Requirement: "an embeddable sandboxed script runtime with a web-style request handler API"

Hosts user scripts in isolated contexts, dispatches HTTP-shaped events to their handlers, and enforces per-context capability limits.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses an HTTP/1.1 request into method, path, headers, body
      - returns error on malformed start line
      # network
    std.http.encode_response
      @ (status: i32, headers: map[string, string], body: bytes) -> bytes
      + serializes an HTTP response to wire format
      # network
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses arbitrary JSON
      - returns error on invalid JSON
      # serialization

script_runtime
  script_runtime.new
    @ () -> runtime_state
    + creates a runtime with no contexts loaded
    # construction
  script_runtime.create_context
    @ (state: runtime_state, name: string, memory_limit_bytes: i64, cpu_budget_ms: i32) -> result[tuple[runtime_state, context_id], string]
    + allocates an isolated execution context with the given quotas
    - returns error when name already exists
    # isolation
  script_runtime.destroy_context
    @ (state: runtime_state, id: context_id) -> runtime_state
    + frees the context and releases its quota
    # isolation
  script_runtime.load_script
    @ (state: runtime_state, id: context_id, source: string) -> result[runtime_state, string]
    + compiles and registers a script inside the context
    - returns error on syntax error
    # loading
  script_runtime.register_handler
    @ (state: runtime_state, id: context_id, route: string, entry: string) -> result[runtime_state, string]
    + maps a route pattern to a named function inside the context
    - returns error when the named function does not exist
    # routing
  script_runtime.grant_capability
    @ (state: runtime_state, id: context_id, capability: string) -> runtime_state
    + enables a named host capability such as "fetch" or "fs_read"
    # security
  script_runtime.revoke_capability
    @ (state: runtime_state, id: context_id, capability: string) -> runtime_state
    + disables a previously granted capability
    # security
  script_runtime.dispatch_request
    @ (state: runtime_state, id: context_id, req: http_request) -> result[http_response, string]
    + routes the request to the matching handler and runs it under the context quotas
    - returns error when no route matches
    - returns error when the script exceeds its CPU or memory budget
    # dispatch
    -> std.http.parse_request
    -> std.time.now_millis
  script_runtime.encode_response
    @ (resp: http_response) -> bytes
    + serializes a handler response for the host transport
    # encoding
    -> std.http.encode_response
  script_runtime.stats
    @ (state: runtime_state, id: context_id) -> optional[context_stats]
    + returns memory used, CPU time consumed, and request count for the context
    - returns none when id is unknown
    # observability
  script_runtime.set_global
    @ (state: runtime_state, id: context_id, key: string, value: json_value) -> runtime_state
    + exposes a host-provided value on the context global object
    # embedding
    -> std.json.parse_value
