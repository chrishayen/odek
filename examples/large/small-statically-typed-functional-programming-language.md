# Requirement: "a small statically typed functional language"

A tiny functional language with type inference. Project runes cover lex, parse, type-check, and evaluate.

std: (all units exist)

fp_lang
  fp_lang.tokenize
    @ (source: string) -> result[list[token], string]
    + produces tokens for identifiers, literals, and keywords
    - returns error on an unrecognized character
    # lexing
  fp_lang.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + builds an AST for let-bindings, lambdas, and applications
    - returns error when a lambda body is missing
    # parsing
  fp_lang.new_type_env
    @ () -> type_env
    + returns an empty type environment
    # environment
  fp_lang.infer
    @ (node: ast_node, env: type_env) -> result[type_scheme, string]
    + returns the inferred type scheme for an expression
    - returns error when unification fails
    # type_inference
  fp_lang.unify
    @ (a: type_scheme, b: type_scheme) -> result[type_scheme, string]
    + returns the most general unifier of two types
    - returns error on an occurs check failure
    # unification
  fp_lang.check
    @ (node: ast_node, env: type_env) -> result[type_env, string]
    + type-checks a top-level declaration and extends the environment
    - returns error when the declared type disagrees with the inferred type
    # type_checking
  fp_lang.new_value_env
    @ () -> value_env
    + returns an empty runtime environment
    # environment
  fp_lang.eval
    @ (node: ast_node, env: value_env) -> result[runtime_value, string]
    + evaluates a type-checked expression
    - returns error when a pattern match is non-exhaustive
    # evaluation
  fp_lang.apply
    @ (fn: runtime_value, arg: runtime_value) -> result[runtime_value, string]
    + applies a closure to an argument
    - returns error when the callee is not a function
    # application
  fp_lang.compile
    @ (source: string) -> result[program, string]
    + tokenizes, parses, and type-checks a source program
    - returns error at the first failing stage
    # pipeline
    -> fp_lang.tokenize
    -> fp_lang.parse
    -> fp_lang.check
  fp_lang.run
    @ (prog: program) -> result[runtime_value, string]
    + evaluates the program's main expression
    - returns error when no main expression is present
    # pipeline
    -> fp_lang.eval
