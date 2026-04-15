# Requirement: "an embeddable dynamic scripting language"

Lex, parse, compile to bytecode, and run on a stack-based virtual machine. Host callers register native functions the script can call.

std: (all units exist)

scriptlang
  scriptlang.tokenize
    fn (source: string) -> result[list[script_token], string]
    + recognizes identifiers, numbers, strings, booleans, operators, and keywords (if, else, while, fn, return)
    - returns error with line and column on unterminated strings
    - returns error on unknown characters
    # lexer
  scriptlang.parse
    fn (tokens: list[script_token]) -> result[script_ast, string]
    + builds an AST of statements and expressions with operator precedence
    - returns error on unexpected tokens with source location
    # parser
  scriptlang.compile
    fn (ast: script_ast) -> result[script_program, string]
    + lowers the AST to a linear bytecode sequence with constant pool
    - returns error on references to undefined names
    # compiler
  scriptlang.new_vm
    fn () -> vm_state
    + creates a fresh virtual machine with empty globals and call stack
    # construction
  scriptlang.register_native
    fn (vm: vm_state, name: string, arity: i32, handler: native_handler) -> vm_state
    + exposes a host-implemented function under the given global name
    # host_binding
  scriptlang.set_global
    fn (vm: vm_state, name: string, value: script_value) -> vm_state
    + assigns a dynamic value to a global name
    # host_binding
  scriptlang.get_global
    fn (vm: vm_state, name: string) -> optional[script_value]
    + retrieves a global value by name
    # host_binding
  scriptlang.run
    fn (vm: vm_state, program: script_program) -> result[script_value, script_error]
    + executes the program and returns the final expression value
    - returns error with source location on runtime type mismatches
    - returns error when the call depth exceeds the configured limit
    # evaluator
  scriptlang.call_function
    fn (vm: vm_state, name: string, args: list[script_value]) -> result[script_value, script_error]
    + invokes a script function by name with runtime arguments
    - returns error when name is not a function
    # evaluator
  scriptlang.eval
    fn (vm: vm_state, source: string) -> result[script_value, script_error]
    + shortcut that tokenizes, parses, compiles, and runs a source string
    # convenience
  scriptlang.format_value
    fn (value: script_value) -> string
    + produces a human-readable representation for debugging
    # display
  scriptlang.stack_trace
    fn (err: script_error) -> list[script_frame]
    + returns the call stack captured at the error site
    # diagnostics
