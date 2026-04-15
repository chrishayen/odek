# Requirement: "a simple optimizing interpreter for a dynamic scripting language"

Parses a tiny expression language, executes it on a stack-based interpreter, and promotes hot functions through an inline cache for property access.

std: (all units exist)

mini_vm
  mini_vm.tokenize
    fn (source: string) -> result[list[tuple[string, string]], string]
    + returns a list of (kind, lexeme) tuples for numbers, identifiers, operators, and keywords
    - returns error on an unterminated string literal
    # lexing
  mini_vm.parse
    fn (tokens: list[tuple[string, string]]) -> result[ast_node, string]
    + parses tokens into an AST rooted at a program node
    - returns error on unexpected end of input
    # parsing
  mini_vm.compile
    fn (ast: ast_node) -> list[i64]
    + lowers the AST to a flat bytecode sequence
    # compilation
  mini_vm.new_vm
    fn () -> vm_state
    + creates a VM with empty stack and globals
    # construction
  mini_vm.execute
    fn (vm: vm_state, bytecode: list[i64]) -> result[tuple[optional[i64], vm_state], string]
    + runs the bytecode and returns the final top-of-stack value if any
    - returns error on stack underflow or type mismatch
    # execution
  mini_vm.inline_cache_lookup
    fn (vm: vm_state, object_shape: i64, property: string) -> optional[i32]
    + returns the cached slot offset for a (shape, property) pair
    - returns none on cache miss
    # optimization
  mini_vm.inline_cache_store
    fn (vm: vm_state, object_shape: i64, property: string, slot: i32) -> vm_state
    + records a slot offset for fast future lookups
    # optimization
  mini_vm.mark_hot
    fn (vm: vm_state, function_id: i64) -> tuple[bool, vm_state]
    + increments a call counter and returns true when the function crosses the optimization threshold
    # optimization
