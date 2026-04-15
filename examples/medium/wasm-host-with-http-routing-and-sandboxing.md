# Requirement: "a host that runs sandboxed WebAssembly modules behind HTTP routes"

A host loads wasm modules, maps them to routes, and invokes exported entry points with request bytes. Module execution lives behind a std wasm runtime primitive.

std
  std.wasm
    std.wasm.instantiate
      fn (module_bytes: bytes) -> result[wasm_instance, string]
      + returns an instance ready to call exported functions
      - returns error when the module is malformed
      # wasm
    std.wasm.call_export
      fn (instance: wasm_instance, name: string, input: bytes) -> result[bytes, string]
      + calls the named exported function with input bytes and returns its output bytes
      - returns error when the export is missing
      - returns error when the guest traps
      # wasm
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from raw request bytes
      - returns error on malformed request line
      # http
    std.http.format_response
      fn (status: i32, headers: map[string, string], body: bytes) -> bytes
      + serializes a status, headers, and body into wire-format bytes
      # http

wasm_host
  wasm_host.new_host
    fn () -> host_state
    + returns an empty host with no modules or routes
    # construction
  wasm_host.load_module
    fn (host: host_state, name: string, module_bytes: bytes) -> result[host_state, string]
    + instantiates the module and registers it under the given name
    - returns error when the module fails to instantiate
    # loading
    -> std.wasm.instantiate
  wasm_host.route
    fn (host: host_state, method: string, path: string, module_name: string, export: string) -> result[host_state, string]
    + maps (method, path) to an exported function of a loaded module
    - returns error when the module name is unknown
    # routing
  wasm_host.invoke
    fn (host: host_state, raw: bytes) -> result[bytes, string]
    + parses the request, finds the matching route, calls the module, and returns a response
    - returns a 404 response when no route matches
    - returns a 500 response when the guest traps
    # dispatch
    -> std.http.parse_request
    -> std.http.format_response
    -> std.wasm.call_export
