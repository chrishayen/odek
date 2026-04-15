# Requirement: "a tree-walking interpreter for a dynamic expression language with variables, functions, and control flow"

The interpreter has three stages: lex, parse, evaluate. Each stage has a narrow entry point. Built-in values are dynamically typed.

std: (all units exist)

interpreter
  interpreter.tokenize
    fn (source: string) -> result[list[token], string]
    + returns a list of tokens for identifiers, numbers, strings, operators, and keywords
    - returns error on an unterminated string literal
    - returns error on an unknown character
    # lexing
  interpreter.parse
    fn (tokens: list[token]) -> result[ast_node, string]
    + returns an AST for a program of statements
    - returns error on unexpected token
    - returns error on unbalanced parentheses or braces
    # parsing
  interpreter.new_env
    fn () -> env_state
    + creates an empty variable environment
    # construction
  interpreter.env_define
    fn (env: env_state, name: string, value: value) -> env_state
    + binds a name to a value in the current scope
    # environment
  interpreter.env_lookup
    fn (env: env_state, name: string) -> result[value, string]
    + returns the value bound to a name
    - returns error when the name is undefined
    # environment
  interpreter.eval
    fn (env: env_state, node: ast_node) -> result[tuple[value, env_state], string]
    + evaluates a program or expression, returning the final value and updated environment
    + supports numeric, string, and boolean literals
    + supports if/else and while control flow
    + supports user-defined functions with lexical closures
    - returns error on type mismatch in arithmetic
    - returns error on division by zero
    - returns error on call to a non-callable value
    # evaluation
  interpreter.run_source
    fn (source: string) -> result[value, string]
    + lexes, parses, and evaluates a source string in a fresh environment
    + returns the value of the last expression
    - returns error at the first failing stage
    # pipeline
  interpreter.call_function
    fn (env: env_state, name: string, args: list[value]) -> result[value, string]
    + looks up a function by name and invokes it with the given arguments
    - returns error when the name is not bound to a function
    - returns error on arity mismatch
    # invocation
  interpreter.register_builtin
    fn (env: env_state, name: string, arity: i32) -> env_state
    + declares a host-provided built-in function in the environment
    ? actual dispatch to host code is resolved by the caller
    # environment
  interpreter.format_value
    fn (v: value) -> string
    + returns a canonical printable string for a value
    + numbers render without trailing zeros, strings render with quotes
    # formatting
