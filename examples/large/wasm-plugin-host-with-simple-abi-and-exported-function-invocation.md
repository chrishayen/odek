# Requirement: "a plug-in host that loads WebAssembly modules, exchanges data with them over a simple host-function ABI, and invokes exported functions"

The host owns module instances and a memory arena; plug-ins read inputs and write outputs through host-provided alloc/free and a small set of host functions.

std
  std.wasm
    std.wasm.load_module
      @ (bytecode: bytes) -> result[wasm_module, string]
      + parses and validates a wasm module
      - returns error on malformed bytecode
      # wasm
    std.wasm.instantiate
      @ (module: wasm_module, imports: map[string, host_fn]) -> result[wasm_instance, string]
      + instantiates a module with the given imported host functions
      - returns error when imports are missing or mistyped
      # wasm
    std.wasm.call
      @ (instance: wasm_instance, fn_name: string, args: list[i64]) -> result[list[i64], string]
      + invokes an exported function
      - returns error when the export does not exist
      - returns error on trap
      # wasm
    std.wasm.read_memory
      @ (instance: wasm_instance, offset: i32, length: i32) -> result[bytes, string]
      + reads a slice of linear memory
      - returns error on out-of-bounds
      # wasm
    std.wasm.write_memory
      @ (instance: wasm_instance, offset: i32, data: bytes) -> result[void, string]
      + writes to linear memory
      - returns error on out-of-bounds
      # wasm
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      - returns error when the file cannot be opened
      # filesystem

plugin_host
  plugin_host.new
    @ () -> host_state
    + returns an empty host
    # construction
  plugin_host.register_host_function
    @ (state: host_state, name: string, impl: host_fn) -> host_state
    + registers a host function callable from wasm
    # registration
  plugin_host.load_from_file
    @ (state: host_state, name: string, path: string) -> result[host_state, string]
    + loads a plug-in from a file and instantiates it with the registered host functions
    - returns error on read or instantiation failure
    # loading
    -> std.fs.read_all
    -> std.wasm.load_module
    -> std.wasm.instantiate
  plugin_host.call_with_bytes
    @ (state: host_state, name: string, export_name: string, input: bytes) -> result[bytes, string]
    + allocates input inside the plug-in's memory, calls the export, reads the output
    - returns error when the plug-in or export is not found
    # invocation
    -> std.wasm.call
    -> std.wasm.read_memory
    -> std.wasm.write_memory
  plugin_host.call_with_string
    @ (state: host_state, name: string, export_name: string, input: string) -> result[string, string]
    + convenience wrapper over call_with_bytes for UTF-8 strings
    # invocation
    -> plugin_host.call_with_bytes
  plugin_host.unload
    @ (state: host_state, name: string) -> host_state
    + drops a plug-in instance and its memory
    # teardown
  plugin_host.list_plugins
    @ (state: host_state) -> list[string]
    + returns the names of loaded plug-ins
    # introspection
