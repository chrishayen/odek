# Requirement: "a dynamic expression evaluator based on s-expressions with an easy extension model"

Parse parenthesized s-expressions, evaluate them against an environment, and let callers register new named functions.

std: (all units exist)

sexpr_eval
  sexpr_eval.parse
    fn (source: string) -> result[expr, string]
    + parses a single s-expression from source
    + supports integer, float, string, symbol, and list atoms
    - returns error on unbalanced parentheses
    - returns error on unterminated string literal
    # parsing
  sexpr_eval.new_env
    fn () -> env_state
    + returns an environment pre-populated with arithmetic and comparison builtins
    ? builtins include + - * / = < > and if
    # construction
  sexpr_eval.bind
    fn (env: env_state, name: string, value: value) -> env_state
    + binds a symbol to a value in the environment
    # environment
  sexpr_eval.register
    fn (env: env_state, name: string, arity: i32, impl_key: string) -> env_state
    + registers an extension function under a symbol with a fixed arity
    # extension
  sexpr_eval.eval
    fn (env: env_state, e: expr) -> result[value, string]
    + evaluates the expression against the environment
    - returns error on unbound symbol
    - returns error on arity mismatch
    - returns error when a non-callable is in function position
    # evaluation
  sexpr_eval.value_int
    fn (n: i64) -> value
    + wraps an integer as a value
    # value
  sexpr_eval.value_string
    fn (s: string) -> value
    + wraps a string as a value
    # value
  sexpr_eval.format_value
    fn (v: value) -> string
    + renders a value as readable text
    # rendering
