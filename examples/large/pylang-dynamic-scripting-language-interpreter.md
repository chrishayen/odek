# Requirement: "an implementation of a dynamic scripting language interpreter"

Full interpreter front-end: source -> tokens -> AST -> bytecode -> VM. Std has no specialized pieces; the interpreter defines everything.

std: (all units exist)

pylang
  pylang.tokenize
    fn (source: string) -> result[list[token], string]
    + returns the token stream for a source file, including indent/dedent tokens
    - returns error on unterminated string literals
    - returns error on mismatched indentation
    # lexing
  pylang.parse
    fn (tokens: list[token]) -> result[ast_node, string]
    + returns a module AST
    - returns error with line and column on syntax errors
    # parsing
  pylang.compile
    fn (ast: ast_node) -> result[code_object, string]
    + returns a code object containing constants, names, and bytecode
    - returns error on references to undeclared free variables
    # compilation
  pylang.vm_new
    fn () -> vm_state
    + returns a fresh interpreter state with empty builtins
    # construction
  pylang.vm_install_builtin
    fn (vm: vm_state, name: string, fn: builtin_fn) -> vm_state
    + returns a new vm with the builtin bound to the given name
    # builtins
  pylang.vm_exec
    fn (vm: vm_state, code: code_object) -> result[tuple[value, vm_state], string]
    + runs the code object and returns its final value plus updated vm
    - returns error on unhandled exceptions
    # execution
  pylang.vm_call
    fn (vm: vm_state, fn: value, args: list[value]) -> result[tuple[value, vm_state], string]
    + invokes a callable value and returns its return value
    - returns error when fn is not callable
    - returns error when arity does not match
    # execution
  pylang.value_to_string
    fn (v: value) -> string
    + returns the printable representation of a value
    # formatting
  pylang.value_equal
    fn (a: value, b: value) -> bool
    + returns true when two values are structurally equal
    # comparison
  pylang.gc_collect
    fn (vm: vm_state) -> vm_state
    + sweeps unreachable values and returns a compacted vm
    # memory
  pylang.eval_string
    fn (vm: vm_state, source: string) -> result[tuple[value, vm_state], string]
    + one-shot: tokenize, parse, compile, and execute a source string
    - returns error on any phase failure with a phase-tagged message
    # convenience
    -> pylang.tokenize
    -> pylang.parse
    -> pylang.compile
    -> pylang.vm_exec
