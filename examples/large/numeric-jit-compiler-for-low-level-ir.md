# Requirement: "a numeric kernel JIT compiler targeting a low-level IR"

Takes a restricted expression AST over numeric arrays, lowers it to an IR, runs a few optimization passes, and emits machine code through a backend.

std
  std.mem
    std.mem.alloc_executable
      fn (size: i64) -> result[bytes, string]
      + allocates an executable memory region of the given size
      - returns error when allocation fails
      # memory
    std.mem.free_executable
      fn (region: bytes) -> void
      + releases a previously allocated executable region
      # memory
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns an FNV-1a 64-bit hash
      # hashing

jit
  jit.parse_kernel
    fn (source: string) -> result[ast_node, string]
    + returns the AST for a restricted arithmetic expression over named inputs
    - returns error on syntax errors
    # frontend
  jit.type_check
    fn (node: ast_node, inputs: map[string, dtype]) -> result[ast_node, string]
    + returns a typed AST when all operand dtypes are consistent
    - returns error when dtypes conflict
    # semantic_analysis
  jit.lower_to_ir
    fn (node: ast_node) -> ir_module
    + lowers a typed AST to a linear IR of SSA instructions
    # lowering
  jit.fold_constants
    fn (module: ir_module) -> ir_module
    + evaluates instructions whose operands are all constant
    # optimization
  jit.eliminate_common_subexpressions
    fn (module: ir_module) -> ir_module
    + replaces duplicate pure instructions with a single definition
    # optimization
  jit.vectorize_inner_loop
    fn (module: ir_module, width: i32) -> ir_module
    + unrolls and vectorizes the innermost loop by the given width
    ? falls back to scalar when the loop body is not vectorizable
    # optimization
  jit.allocate_registers
    fn (module: ir_module, available: i32) -> ir_module
    + assigns virtual registers to physical registers, inserting spills if needed
    # backend
  jit.emit_machine_code
    fn (module: ir_module) -> bytes
    + returns the encoded machine code for the module
    # backend
  jit.cache_lookup
    fn (state: jit_cache, key: string) -> optional[compiled_kernel]
    + returns a previously compiled kernel for the cache key
    # caching
    -> std.hash.fnv64
  jit.compile
    fn (state: jit_cache, source: string, inputs: map[string, dtype]) -> result[compiled_kernel, string]
    + returns a callable kernel from source, reusing cached code when possible
    - returns error on parse, type, or codegen failure
    # compilation
    -> std.mem.alloc_executable
  jit.invoke
    fn (kernel: compiled_kernel, args: list[buffer]) -> result[void, string]
    + executes a compiled kernel over the given input/output buffers
    - returns error when argument shapes do not match the kernel signature
    # execution
  jit.release
    fn (kernel: compiled_kernel) -> void
    + frees the machine code region held by the kernel
    # lifecycle
    -> std.mem.free_executable
