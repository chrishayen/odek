# Requirement: "a dynamic language interpreter"

A tree-walking interpreter for a dynamic scripting language: lexing, parsing, evaluating with an environment, and supporting first-class functions.

std
  std.io
    std.io.read_file
      @ (path: string) -> result[string, string]
      + returns file contents as a string
      - returns error when path does not exist
      # io
  std.strings
    std.strings.is_digit
      @ (c: string) -> bool
      + returns true for characters '0'..'9'
      - returns false for letters and punctuation
      # text
    std.strings.is_alpha
      @ (c: string) -> bool
      + returns true for ascii letters and underscore
      - returns false for digits and punctuation
      # text
  std.collections
    std.collections.hashmap_new
      @ () -> map[string, value]
      + returns an empty map
      # collections
    std.collections.hashmap_get
      @ (m: map[string, value], key: string) -> optional[value]
      + returns the value associated with key
      - returns none when key missing
      # collections

interpreter
  interpreter.tokenize
    @ (source: string) -> result[list[token], string]
    + produces tokens for numbers, identifiers, strings, operators, and keywords
    - returns error on unterminated string literal
    - returns error on unknown character
    # lexing
    -> std.strings.is_digit
    -> std.strings.is_alpha
  interpreter.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + parses tokens into an AST of statements and expressions
    + supports function definitions, if/else, while, return, and binary ops
    - returns error on unexpected token
    - returns error on unclosed parenthesis or brace
    # parsing
  interpreter.new_environment
    @ (parent: optional[environment]) -> environment
    + creates a new scope optionally nested inside a parent scope
    # scoping
    -> std.collections.hashmap_new
  interpreter.env_define
    @ (env: environment, name: string, val: value) -> environment
    + binds a name to a value in the current scope
    # scoping
  interpreter.env_lookup
    @ (env: environment, name: string) -> result[value, string]
    + walks the scope chain and returns the bound value
    - returns error when name is undefined
    # scoping
    -> std.collections.hashmap_get
  interpreter.eval_expression
    @ (node: ast_node, env: environment) -> result[value, string]
    + evaluates literals, variables, binary ops, and calls
    + returns a first-class function value for function literals
    - returns error on type mismatch in binary op
    - returns error on call to non-callable value
    # evaluation
  interpreter.eval_statement
    @ (node: ast_node, env: environment) -> result[environment, string]
    + executes a statement and returns the updated environment
    + handles if, while, return, and expression statements
    - returns error on runtime failure in nested expression
    # evaluation
  interpreter.call_function
    @ (func: value, args: list[value], env: environment) -> result[value, string]
    + binds arguments to parameters in a new child scope and evaluates the body
    - returns error on arity mismatch
    # function_calls
  interpreter.run_program
    @ (source: string) -> result[value, string]
    + runs a full program and returns the final expression value
    - returns error at the first compile or runtime failure
    # entrypoint
