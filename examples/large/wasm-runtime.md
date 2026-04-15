# Requirement: "a lightweight runtime for webassembly"

A minimal interpreter that parses a module binary, instantiates it, and executes exported functions on a value stack.

std
  std.bytes
    std.bytes.read_u32_le
      fn (data: bytes, offset: i32) -> result[tuple[u32, i32], string]
      + reads a little-endian u32 and returns the value with the new offset
      - returns error when there are fewer than 4 bytes available
      # binary_reading
    std.bytes.read_leb128_u32
      fn (data: bytes, offset: i32) -> result[tuple[u32, i32], string]
      + decodes an unsigned LEB128 integer up to 32 bits
      - returns error on truncated input or overflow
      # binary_reading
    std.bytes.read_leb128_i32
      fn (data: bytes, offset: i32) -> result[tuple[i32, i32], string]
      + decodes a signed LEB128 integer up to 32 bits
      - returns error on truncated input or overflow
      # binary_reading

wasm
  wasm.parse_module
    fn (binary: bytes) -> result[module, string]
    + decodes the magic number, version, and every known section
    - returns error when magic is not "\0asm" or version is unsupported
    - returns error on malformed section headers
    # parsing
    -> std.bytes.read_u32_le
    -> std.bytes.read_leb128_u32
  wasm.parse_type_section
    fn (data: bytes) -> result[list[func_type], string]
    + parses the type section into a list of function signatures
    - returns error on unknown value types
    # parsing
    -> std.bytes.read_leb128_u32
  wasm.parse_function_section
    fn (data: bytes) -> result[list[u32], string]
    + parses the function section into a list of type indices
    # parsing
    -> std.bytes.read_leb128_u32
  wasm.parse_code_section
    fn (data: bytes) -> result[list[func_body], string]
    + parses locals and instruction bytes for each function
    - returns error when declared body size does not match consumed bytes
    # parsing
    -> std.bytes.read_leb128_u32
  wasm.parse_export_section
    fn (data: bytes) -> result[map[string, export_ref], string]
    + parses name-to-reference pairs for exported functions, tables, memories, globals
    # parsing
    -> std.bytes.read_leb128_u32
  wasm.instantiate
    fn (m: module) -> result[instance, string]
    + allocates memory pages and globals declared by the module
    - returns error when a required import is unsatisfied
    # instantiation
  wasm.invoke
    fn (inst: instance, export_name: string, args: list[value]) -> result[list[value], string]
    + looks up the named export and runs its body on a fresh value stack
    - returns error when the export is missing or not a function
    - returns error on arity mismatch
    # execution
  wasm.step
    fn (inst: instance, frame: call_frame) -> result[call_frame, string]
    + executes one instruction and returns the updated frame
    - returns error on unknown opcode or stack underflow
    -> std.bytes.read_leb128_i32
    # execution
  wasm.decode_value_type
    fn (byte_val: u8) -> result[value_type, string]
    + maps a byte tag to i32, i64, f32, or f64
    - returns error on unknown tag
    # parsing
