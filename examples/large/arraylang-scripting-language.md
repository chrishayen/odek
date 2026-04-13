# Requirement: "an embeddable array-oriented scripting language"

Values are arrays; operators apply element-wise. Lex, parse, evaluate, and allow host bindings.

std: (all units exist)

arraylang
  arraylang.tokenize
    @ (source: string) -> result[list[array_token], string]
    + recognizes numbers, strings, identifiers, primitives, and brackets
    - returns error with position on unterminated strings
    # lexer
  arraylang.parse
    @ (tokens: list[array_token]) -> result[array_ast, string]
    + builds an AST honoring right-to-left associativity typical for array languages
    - returns error on unbalanced brackets
    # parser
  arraylang.new_env
    @ () -> array_env
    + creates an empty environment with the built-in primitives preloaded
    # construction
  arraylang.bind
    @ (env: array_env, name: string, value: array_value) -> array_env
    + assigns a value to an identifier in the environment
    # binding
  arraylang.register_native
    @ (env: array_env, name: string, handler: native_handler) -> array_env
    + exposes a host function as a named verb
    # host_binding
  arraylang.eval
    @ (env: array_env, ast: array_ast) -> result[tuple[array_value, array_env], string]
    + evaluates the AST and returns the final value and updated environment
    - returns error on arity or shape mismatches
    # evaluator
  arraylang.run
    @ (env: array_env, source: string) -> result[tuple[array_value, array_env], string]
    + tokenizes, parses, and evaluates a source string
    # convenience
  arraylang.apply_unary
    @ (verb: string, arg: array_value) -> result[array_value, string]
    + applies a primitive unary verb element-wise
    - returns error on unknown verb
    # primitives
  arraylang.apply_binary
    @ (verb: string, left: array_value, right: array_value) -> result[array_value, string]
    + applies a primitive binary verb element-wise with broadcasting
    - returns error on incompatible shapes
    # primitives
  arraylang.format_value
    @ (value: array_value) -> string
    + produces a printed representation aligned by columns
    # display
  arraylang.reshape
    @ (shape: list[i32], value: array_value) -> result[array_value, string]
    + rearranges an array into the given shape
    - returns error when the element count does not match
    # shape
