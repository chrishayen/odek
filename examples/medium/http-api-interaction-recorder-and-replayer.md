# Requirement: "a proxy that records and replays http api interactions"

Runs in record mode to capture upstream responses, or simulate mode to serve matching captured responses. Middleware hooks can mutate requests or responses.

std
  std.http
    std.http.forward
      @ (url: string, method: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + forwards the request upstream and returns the response
      - returns error on connection failure
      # http
  std.io
    std.io.read_file
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when the path does not exist
      # filesystem
    std.io.write_file
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to path, creating or truncating
      # filesystem
  std.json
    std.json.parse
      @ (raw: bytes) -> result[json_value, string]
      + parses raw bytes as a json document
      - returns error on invalid json
      # serialization
    std.json.encode
      @ (value: json_value) -> bytes
      + encodes a json value to compact bytes
      # serialization

record_proxy
  record_proxy.new
    @ (mode: string) -> proxy_state
    + returns a proxy in "record" or "simulate" mode with no captures
    - returns an empty state for an unknown mode string
    # construction
  record_proxy.handle
    @ (state: proxy_state, req: http_request) -> result[tuple[http_response, proxy_state], string]
    + in record mode, forwards upstream and stores the pair keyed by method+path+body-hash
    + in simulate mode, returns the stored response for a matching key
    - in simulate mode, returns error when no capture matches
    # dispatch
    -> std.http.forward
  record_proxy.save
    @ (state: proxy_state, path: string) -> result[void, string]
    + serializes captured pairs to a json file
    # persistence
    -> std.json.encode
    -> std.io.write_file
  record_proxy.load
    @ (path: string) -> result[proxy_state, string]
    + reads captures from a json file into a new proxy state
    - returns error on invalid json
    # persistence
    -> std.io.read_file
    -> std.json.parse
  record_proxy.add_request_hook
    @ (state: proxy_state, hook: request_hook) -> proxy_state
    + appends a middleware invoked on every incoming request before matching
    # middleware
  record_proxy.add_response_hook
    @ (state: proxy_state, hook: response_hook) -> proxy_state
    + appends a middleware invoked on every response before returning
    # middleware
