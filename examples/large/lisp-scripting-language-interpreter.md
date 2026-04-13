# Requirement: "a lisp-like scripting language interpreter"

Source text goes through a reader, into an AST, then into a tree-walking evaluator over an environment. Builtin primitives and garbage collection are part of the core.

std: (all units exist)

gamelisp
  gamelisp.read
    @ (src: string) -> result[list[form], string]
    + parses source into a list of top-level forms
    - returns error on unbalanced parentheses or invalid tokens
    # reader
  gamelisp.read_form
    @ (tokens: list[token], pos: i32) -> result[tuple[form, i32], string]
    + parses one form starting at pos and returns (form, next_pos)
    # reader
  gamelisp.tokenize
    @ (src: string) -> result[list[token], string]
    + produces tokens for parens, symbols, numbers, strings, quotes
    - returns error on unterminated string literal
    # lexer
  gamelisp.new_env
    @ () -> env_state
    + returns an environment pre-populated with core builtins
    # construction
  gamelisp.define
    @ (env: env_state, name: string, value: value) -> env_state
    + binds name to value in the current scope
    # environment
  gamelisp.lookup
    @ (env: env_state, name: string) -> result[value, string]
    + resolves name walking enclosing scopes
    - returns error when name is unbound
    # environment
  gamelisp.eval
    @ (env: env_state, f: form) -> result[value, string]
    + evaluates one form, handling special forms (quote, if, let, fn, def, set!, do)
    - returns error on unknown symbol or arity mismatch
    # evaluator
  gamelisp.apply
    @ (env: env_state, callee: value, args: list[value]) -> result[value, string]
    + applies a function value to arguments, dispatching to builtins or closures
    - returns error when callee is not callable
    # evaluator
  gamelisp.make_closure
    @ (env: env_state, params: list[string], body: list[form]) -> value
    + captures env and returns a closure value
    # evaluator
  gamelisp.register_builtin
    @ (env: env_state, name: string, arity: i32, impl_id: i32) -> env_state
    + registers a native function callable from scripts
    # interop
  gamelisp.gc_collect
    @ (env: env_state) -> env_state
    + sweeps unreachable values from the managed heap
    ? uses mark-and-sweep over live roots
    # memory
  gamelisp.run
    @ (env: env_state, src: string) -> result[value, string]
    + reads and evaluates src, returning the last form's value
    - returns error when read or eval fails
    # entry_point
