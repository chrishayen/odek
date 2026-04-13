# Requirement: "a static type checker for a dynamically typed scripting language"

A lexer, parser, and type inference engine over a small expression language with functions, primitives, and lists. The project layer runs analyses; std holds general-purpose parsing helpers.

std
  std.io
    std.io.read_file
      @ (path: string) -> result[string, string]
      + reads a file as UTF-8 text
      - returns error when the path does not exist
      # filesystem

type_checker
  type_checker.tokenize
    @ (source: string) -> result[list[token], string]
    + returns tokens in source order
    - returns error with position on an unrecognized character
    # lexing
  type_checker.parse
    @ (tokens: list[token]) -> result[ast_node, string]
    + returns the root AST node for the source
    - returns error with position on a syntax error
    # parsing
  type_checker.new_environment
    @ () -> type_environment
    + returns an empty environment with built-in types preloaded
    # construction
  type_checker.bind
    @ (env: type_environment, name: string, ty: inferred_type) -> type_environment
    + returns a new environment with the binding added
    # environment
  type_checker.infer
    @ (env: type_environment, node: ast_node) -> result[inferred_type, type_error]
    + returns the inferred type of an expression
    - returns error when a function is applied to the wrong number of arguments
    - returns error when a binary operator receives incompatible operand types
    - returns error when a free variable is referenced
    # inference
  type_checker.unify
    @ (a: inferred_type, b: inferred_type) -> result[inferred_type, type_error]
    + returns the most general type that satisfies both inputs
    - returns error when the types have no common instantiation
    # unification
  type_checker.check_program
    @ (source: string) -> result[list[inferred_type], list[type_error]]
    + returns the inferred type of every top-level binding when the program is well-typed
    - returns the collected errors when any checks fail
    # orchestration
  type_checker.check_file
    @ (path: string) -> result[list[inferred_type], list[type_error]]
    + checks a source file from disk
    - returns error when the file cannot be read
    # orchestration
    -> std.io.read_file
  type_checker.format_error
    @ (err: type_error) -> string
    + renders a type error as a human-readable diagnostic with source position
    # diagnostics
  type_checker.fresh_variable
    @ (env: type_environment) -> tuple[type_environment, inferred_type]
    + returns the environment with a new unique type variable allocated
    # inference
  type_checker.substitute
    @ (ty: inferred_type, bindings: map[string, inferred_type]) -> inferred_type
    + returns the type with every variable replaced by its binding
    # inference
