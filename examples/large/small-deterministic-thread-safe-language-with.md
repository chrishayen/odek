# Requirement: "a small deterministic embedded scripting language"

A tiny interpreter for a Python-syntax scripting language. Evaluation is deterministic: no clocks, no randomness, no I/O. The project layer covers lex, parse, and evaluate against an environment.

std
  std.strings
    std.strings.split_lines
      @ (s: string) -> list[string]
      + splits on newline preserving order
      # strings
  std.collections
    std.collections.map_get
      @ (m: map[string, i64], key: string) -> optional[i64]
      + returns the value when present
      - returns none when the key is absent
      # collections

script_lang
  script_lang.tokenize
    @ (source: string) -> result[list[token], string]
    + produces tokens for identifiers, integers, operators, and keywords
    - returns error on an unterminated string literal
    # lexing
    -> std.strings.split_lines
  script_lang.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + builds an AST for expressions and statements
    - returns error when parentheses are unbalanced
    # parsing
  script_lang.new_env
    @ () -> env_state
    + returns an empty variable environment
    # environment
  script_lang.bind
    @ (env: env_state, name: string, value: i64) -> env_state
    + returns a new environment with the binding added
    ? environments are immutable to keep evaluation deterministic
    # environment
  script_lang.eval_expr
    @ (node: ast_node, env: env_state) -> result[i64, string]
    + evaluates arithmetic and boolean expressions
    - returns error when a referenced variable is unbound
    # evaluation
    -> std.collections.map_get
  script_lang.eval_block
    @ (nodes: list[ast_node], env: env_state) -> result[env_state, string]
    + executes statements sequentially, threading the environment
    - returns error on the first failing statement
    # evaluation
  script_lang.call_function
    @ (env: env_state, name: string, args: list[i64]) -> result[i64, string]
    + invokes a user-defined function and returns its result
    - returns error when arity does not match
    # function_call
  script_lang.run
    @ (source: string, env: env_state) -> result[env_state, string]
    + tokenizes, parses, and evaluates a program
    - returns error at the first stage that fails
    # pipeline
    -> script_lang.tokenize
    -> script_lang.parse
    -> script_lang.eval_block
