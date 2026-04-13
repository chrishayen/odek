# Requirement: "an embeddable general-purpose scripting language"

A conventional dynamic language: lexer, recursive-descent parser, tree-walking evaluator, and host interop.

std: (all units exist)

scriptcore
  scriptcore.tokenize
    @ (source: string) -> result[list[script_token], string]
    + recognizes identifiers, keywords (if, else, while, for, fn, return), numbers, strings, operators
    - returns error with line and column on unterminated strings
    # lexer
  scriptcore.parse
    @ (tokens: list[script_token]) -> result[script_ast, string]
    + builds a statement-level AST with expression precedence
    - returns error on unexpected tokens
    # parser
  scriptcore.new_env
    @ () -> script_env
    + creates a root environment with the built-in functions preloaded
    # construction
  scriptcore.register_native
    @ (env: script_env, name: string, arity: i32, handler: native_handler) -> script_env
    + exposes a host function under a name
    # host_binding
  scriptcore.set
    @ (env: script_env, name: string, value: script_value) -> script_env
    + assigns a value to a variable in the current scope
    # binding
  scriptcore.get
    @ (env: script_env, name: string) -> optional[script_value]
    + resolves a variable by walking scopes outward
    # binding
  scriptcore.eval
    @ (env: script_env, ast: script_ast) -> result[tuple[script_value, script_env], script_error]
    + executes the AST, returning the last expression value
    - returns error on runtime type mismatches
    - returns error on call depth overflow
    # evaluator
  scriptcore.call
    @ (env: script_env, name: string, args: list[script_value]) -> result[script_value, script_error]
    + invokes a script or native function by name
    - returns error when name does not resolve to a callable
    # evaluator
  scriptcore.run
    @ (env: script_env, source: string) -> result[tuple[script_value, script_env], script_error]
    + tokenizes, parses, and evaluates a source string
    # convenience
  scriptcore.format_value
    @ (value: script_value) -> string
    + produces a debug representation of a runtime value
    # display
  scriptcore.stack_trace
    @ (err: script_error) -> list[script_frame]
    + returns the call chain captured at the error
    # diagnostics
