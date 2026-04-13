# Requirement: "a small language for GPGPU numerical computing"

A frontend that lexes, parses, and type-checks kernel source, then lowers it to an intermediate representation suitable for GPU backends. Execution happens through a pluggable device interface.

std
  std.io
    std.io.read_all
      @ (path: string) -> result[string, string]
      + reads a file into a string
      - returns error when the file does not exist
      # io
  std.collections
    std.collections.map_get
      @ (m: map[string,string], key: string) -> optional[string]
      + returns the value associated with key if present
      # collections

gpulang
  gpulang.tokenize
    @ (source: string) -> result[list[token], string]
    + produces a token stream from source text
    - returns error on unterminated string or invalid character
    # lexing
  gpulang.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + builds an abstract syntax tree from a token stream
    - returns error on unexpected token
    # parsing
  gpulang.typecheck
    @ (ast: ast_node) -> result[typed_ast, string]
    + annotates every expression with a numeric type
    - returns error when a scalar and a vector are combined without broadcast rules
    # semantic_analysis
  gpulang.lower_to_ir
    @ (typed: typed_ast) -> ir_module
    + converts a typed AST into a flat IR with explicit memory operations
    # ir_lowering
  gpulang.allocate_buffer
    @ (device: gpu_device, size_bytes: i64) -> result[buffer_handle, string]
    + reserves device memory and returns a handle
    - returns error when the device reports insufficient memory
    # memory
  gpulang.upload
    @ (device: gpu_device, buf: buffer_handle, data: bytes) -> result[void, string]
    + copies host bytes into a device buffer
    - returns error when data size exceeds buffer capacity
    # memory
  gpulang.compile_kernel
    @ (device: gpu_device, ir: ir_module, entry_point: string) -> result[kernel_handle, string]
    + compiles an IR module for the target device
    - returns error when the entry point is not defined
    # compilation
  gpulang.launch
    @ (device: gpu_device, kernel: kernel_handle, grid: list[i32], args: list[buffer_handle]) -> result[void, string]
    + dispatches a kernel over the given grid dimensions
    - returns error on invalid grid shape
    # execution
  gpulang.download
    @ (device: gpu_device, buf: buffer_handle) -> result[bytes, string]
    + copies device buffer contents back to host memory
    # memory
  gpulang.compile_source
    @ (source: string) -> result[ir_module, string]
    + full pipeline from source string to IR
    - returns error at any stage with a diagnostic message
    # pipeline
    -> std.io.read_all
