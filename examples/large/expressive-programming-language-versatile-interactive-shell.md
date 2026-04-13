# Requirement: "a programming language and interactive shell library"

The library exposes a small expression language: lex, parse, evaluate. An interactive-shell layer maintains variable bindings and history across lines. Executing external programs is delegated to a std process primitive.

std
  std.process
    std.process.spawn_and_wait
      @ (program: string, args: list[string], stdin: bytes) -> result[process_result, string]
      + runs the program to completion and returns stdout, stderr, and exit code
      - returns error when the program cannot be launched
      # process
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of a regular file
      # filesystem

shlang
  shlang.tokenize
    @ (source: string) -> result[list[token], string]
    + returns tokens for identifiers, numbers, strings, operators, and pipes
    - returns error on unterminated string literal
    # lexing
  shlang.parse_expression
    @ (tokens: list[token]) -> result[expr, string]
    + returns an AST for a single expression
    - returns error on unexpected token
    - returns error on unmatched brackets
    # parsing
  shlang.parse_pipeline
    @ (tokens: list[token]) -> result[pipeline, string]
    + returns a pipeline of commands separated by pipes
    # parsing
  shlang.new_env
    @ () -> env_state
    + creates an env with an empty binding table and history
    # state
  shlang.bind
    @ (env: env_state, name: string, value: value) -> env_state
    + records a variable binding
    # state
  shlang.lookup
    @ (env: env_state, name: string) -> optional[value]
    + returns the value bound to name
    # state
  shlang.eval_expression
    @ (env: env_state, e: expr) -> result[value, string]
    + evaluates an expression against the environment
    - returns error on unbound name
    - returns error on type mismatch
    # evaluation
  shlang.run_pipeline
    @ (env: env_state, p: pipeline) -> result[value, string]
    + runs a pipeline, executing external commands and piping stdout forward
    - returns error when any stage fails to launch
    # execution
    -> std.process.spawn_and_wait
  shlang.eval_line
    @ (env: env_state, line: string) -> tuple[result[value, string], env_state]
    + tokenizes, parses, evaluates a single input line, and appends it to history
    # repl
  shlang.history
    @ (env: env_state) -> list[string]
    + returns the history of lines evaluated against this environment
    # state
  shlang.load_script
    @ (env: env_state, path: string) -> tuple[result[value, string], env_state]
    + reads a script file and evaluates it as a sequence of lines
    # scripting
    -> std.fs.read_all
