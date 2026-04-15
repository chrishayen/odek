# Requirement: "a lean and efficient implementation of a dynamic programming language"

A compact interpreter pipeline: lexer, parser, bytecode compiler, and stack-based virtual machine.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a source file into a string
      - returns error when the file does not exist
      # filesystem

minipy
  minipy.tokenize
    fn (source: string) -> result[list[token], string]
    + produces a token stream from source text with indent/dedent tokens
    - returns error on inconsistent indentation
    # lexing
  minipy.parse
    fn (tokens: list[token]) -> result[ast_node, string]
    + builds an AST from the token stream
    - returns error on unexpected token
    # parsing
  minipy.compile
    fn (ast: ast_node) -> bytecode_module
    + lowers an AST into a bytecode module with a constant pool
    # compilation
  minipy.new_vm
    fn () -> vm_state
    + creates an empty VM with an operand stack and empty globals
    # vm
  minipy.run_module
    fn (vm: vm_state, module: bytecode_module) -> result[value, string]
    + executes a bytecode module and returns the top-of-stack value
    - returns error on stack underflow or undefined name
    # execution
  minipy.call_function
    fn (vm: vm_state, name: string, args: list[value]) -> result[value, string]
    + invokes a named function in the current globals
    - returns error when the name is not callable
    # execution
  minipy.define_builtin
    fn (vm: vm_state, name: string, fn: builtin_handle) -> vm_state
    + registers a host-provided function under a name
    # extension
  minipy.load_module
    fn (vm: vm_state, path: string) -> result[bytecode_module, string]
    + reads, compiles, and returns a module from a source file
    # loader
    -> std.fs.read_all
  minipy.collect_garbage
    fn (vm: vm_state) -> vm_state
    + performs a mark-and-sweep pass over reachable values
    ? simple collector; assumes no weak references
    # memory
  minipy.reset_vm
    fn (vm: vm_state) -> vm_state
    + clears the operand stack and local frames, retaining globals
    # state
