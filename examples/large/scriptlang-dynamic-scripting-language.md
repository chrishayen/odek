# Requirement: "an embeddable dynamic scripting language"

Lex, parse, compile to bytecode, and run on a stack-based virtual machine. Host callers register native functions the script can call.

std: (all units exist)

scriptlang
  scriptlang.tokenize
    @ (source: string) -> result[list[script_token], string]
    + recognizes identifiers, numbers, strings, booleans, operators, and keywords (if, else, while, fn, return)
    - returns error with line and column on unterminated strings
    - returns error on unknown characters
    # lexer
  scriptlang.parse
    @ (tokens: list[script_token]) -> result[script_ast, string]
    + builds an AST of statements and expressions with operator precedence
    - returns error on unexpected tokens with source location
    # parser
  scriptlang.compile
    @ (ast: script_ast) -> result[script_program, string]
    + lowers the AST to a linear bytecode sequence with constant pool
    - returns error on references to undefined names
    # compiler
  scriptlang.new_vm
    @ () -> vm_state
    + creates a fresh virtual machine with empty globals and call stack
    # construction
  scriptlang.register_native
    @ (vm: vm_state, name: string, arity: i32, handler: native_handler) -> vm_state
    + exposes a host-implemented function under the given global name
    # host_binding
  scriptlang.set_global
    @ (vm: vm_state, name: string, value: script_value) -> vm_state
    + assigns a dynamic value to a global name
    # host_binding
  scriptlang.get_global
    @ (vm: vm_state, name: string) -> optional[script_value]
    + retrieves a global value by name
    # host_binding
  scriptlang.run
    @ (vm: vm_state, program: script_program) -> result[script_value, script_error]
    + executes the program and returns the final expression value
    - returns error with source location on runtime type mismatches
    - returns error when the call depth exceeds the configured limit
    # evaluator
  scriptlang.call_function
    @ (vm: vm_state, name: string, args: list[script_value]) -> result[script_value, script_error]
    + invokes a script function by name with runtime arguments
    - returns error when name is not a function
    # evaluator
  scriptlang.eval
    @ (vm: vm_state, source: string) -> result[script_value, script_error]
    + shortcut that tokenizes, parses, compiles, and runs a source string
    # convenience
  scriptlang.format_value
    @ (value: script_value) -> string
    + produces a human-readable representation for debugging
    # display
  scriptlang.stack_trace
    @ (err: script_error) -> list[script_frame]
    + returns the call stack captured at the error site
    # diagnostics
