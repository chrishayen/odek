# Requirement: "a library facilitating high-level interactions between a wasm guest and a host environment"

Bridges typed host functions to a wasm module by marshalling arguments through linear memory.

std
  std.bytes
    std.bytes.read_u32_le
      @ (data: bytes, offset: i32) -> result[u32, string]
      + returns the little-endian u32 at offset
      - returns error on out-of-range offset
      # bytes
    std.bytes.write_u32_le
      @ (data: bytes, offset: i32, value: u32) -> result[bytes, string]
      + writes value little-endian at offset
      - returns error on out-of-range offset
      # bytes
    std.bytes.slice
      @ (data: bytes, start: i32, end: i32) -> bytes
      + returns the subrange [start, end)
      # bytes

wasm_bridge
  wasm_bridge.new
    @ (memory_size: i32) -> bridge_state
    + creates a bridge with a linear memory of the given size
    # construction
  wasm_bridge.alloc
    @ (state: bridge_state, size: i32) -> result[tuple[i32, bridge_state], string]
    + reserves a region of size bytes and returns its offset
    - returns error when memory is exhausted
    # memory
  wasm_bridge.write_string
    @ (state: bridge_state, text: string) -> result[tuple[i32, i32, bridge_state], string]
    + copies text bytes into memory and returns (offset, length)
    # marshalling
    -> std.bytes.write_u32_le
  wasm_bridge.read_string
    @ (state: bridge_state, offset: i32, length: i32) -> result[string, string]
    + returns the UTF-8 string starting at offset
    - returns error when bytes are not valid UTF-8
    # marshalling
    -> std.bytes.slice
  wasm_bridge.register_host_func
    @ (state: bridge_state, name: string, signature: string, func_id: i64) -> result[bridge_state, string]
    + makes a host function callable from the guest
    - returns error when a function with the same name already exists
    # host_binding
  wasm_bridge.invoke_guest
    @ (state: bridge_state, export_name: string, args: list[wasm_value]) -> result[tuple[wasm_value, bridge_state], string]
    + marshals args, dispatches to an export, and returns the typed result
    - returns error when export is unknown
    - returns error on argument/signature mismatch
    # invocation
    -> std.bytes.write_u32_le
    -> std.bytes.read_u32_le
  wasm_bridge.free
    @ (state: bridge_state, offset: i32) -> result[bridge_state, string]
    + returns a previously allocated region to the free list
    - returns error when offset is not a live allocation
    # memory
