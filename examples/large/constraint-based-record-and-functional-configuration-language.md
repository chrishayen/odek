# Requirement: "a constraint-based record and functional configuration language"

A small configuration language with records, field constraints, and pure function application. The project layer covers lexing, parsing, type-checking against constraints, and evaluation.

std
  std.text
    std.text.scan_tokens
      fn (source: string) -> result[list[token], string]
      + splits source into tokens with line and column positions
      - returns error on an unterminated string literal
      # lexing
  std.collections
    std.collections.map_merge
      fn (a: map[string, value], b: map[string, value]) -> map[string, value]
      + returns a new map with b's entries overriding a's
      # collections

config_lang
  config_lang.parse
    fn (source: string) -> result[ast_node, string]
    + parses source into an AST of records, field definitions, and expressions
    - returns error with line/column on a syntax error
    # parsing
    -> std.text.scan_tokens
  config_lang.build_schema
    fn (ast: ast_node) -> result[schema, string]
    + extracts record type definitions and their field constraints from the AST
    - returns error when two record types share a name
    # schema
  config_lang.check_constraints
    fn (schema: schema, record_name: string, fields: map[string, value]) -> result[void, list[string]]
    + returns ok when every field satisfies its declared constraint
    - returns a list of violation messages, one per failing field
    # validation
  config_lang.apply_function
    fn (env: eval_env, name: string, args: list[value]) -> result[value, string]
    + applies a user-defined pure function from the environment to the given arguments
    - returns error when arity does not match
    - returns error when a referenced name is not bound
    # evaluation
  config_lang.evaluate
    fn (ast: ast_node, env: eval_env) -> result[value, string]
    + evaluates an expression AST under an environment of bindings
    - returns error on unbound identifiers or type mismatches
    # evaluation
  config_lang.build_env
    fn (bindings: map[string, value], parent: optional[eval_env]) -> eval_env
    + constructs a lexical environment with optional parent scope
    # environment
    -> std.collections.map_merge
  config_lang.render_record
    fn (value: value) -> result[string, string]
    + serializes a fully evaluated record value back to the surface syntax
    - returns error when the value contains an unevaluated expression
    # rendering
