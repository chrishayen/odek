# Requirement: "a WebAssembly sandboxing runtime for executing untrusted code"

Loads a WASM module, enforces resource limits, and invokes exported functions with isolated memory.

std
  std.wasm
    std.wasm.parse_module
      @ (bytes: bytes) -> result[wasm_module, string]
      + parses and validates a WASM binary
      - returns error on invalid magic or section structure
      # wasm
    std.wasm.instantiate
      @ (module: wasm_module, imports: map[string, host_fn]) -> result[wasm_instance, string]
      + instantiates a module with resolved host imports
      - returns error when required imports are missing
      # wasm
    std.wasm.call
      @ (instance: wasm_instance, name: string, args: list[i64]) -> result[list[i64], string]
      + calls an exported function with integer arguments
      - returns error when the export does not exist or signatures mismatch
      # wasm
    std.wasm.read_memory
      @ (instance: wasm_instance, offset: i32, length: i32) -> result[bytes, string]
      + reads bytes from the instance's linear memory
      - returns error on out-of-bounds access
      # wasm

sandbox
  sandbox.new_limits
    @ (max_memory_pages: i32, max_fuel: i64, max_stack: i32) -> limits
    + creates a resource limit bundle
    # configuration
  sandbox.load
    @ (code: bytes, limits: limits) -> result[sandbox_state, string]
    + parses and instantiates the module with the given limits
    - returns error when the module exceeds declared memory limits
    # loading
    -> std.wasm.parse_module
    -> std.wasm.instantiate
  sandbox.register_host_fn
    @ (state: sandbox_state, name: string, fn: host_fn) -> sandbox_state
    + registers a host function by import name before instantiation
    # host_bindings
  sandbox.invoke
    @ (state: sandbox_state, name: string, args: list[i64]) -> result[list[i64], string]
    + calls an exported function while enforcing the fuel budget
    - returns error when fuel is exhausted
    - returns error when the guest traps
    # execution
    -> std.wasm.call
  sandbox.consume_fuel
    @ (state: sandbox_state, amount: i64) -> result[sandbox_state, string]
    + decrements remaining fuel by the given amount
    - returns error when amount exceeds remaining fuel
    # metering
  sandbox.remaining_fuel
    @ (state: sandbox_state) -> i64
    + returns current fuel budget
    # introspection
  sandbox.copy_out
    @ (state: sandbox_state, offset: i32, length: i32) -> result[bytes, string]
    + copies a region of guest memory out to the host
    - returns error on out-of-bounds
    # memory
    -> std.wasm.read_memory
