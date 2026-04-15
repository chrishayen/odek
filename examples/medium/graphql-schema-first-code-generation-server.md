# Requirement: "a graphql server library driven by schema-first code generation"

The runtime side: a schema registry, resolver dispatch, and query execution. Code generation is a separate offline concern; this library exposes the executor that generated resolvers plug into.

std
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

graphql
  graphql.new_schema
    fn () -> schema_state
    + returns an empty schema registry
    # construction
  graphql.register_type
    fn (state: schema_state, type_name: string, field_names: list[string]) -> schema_state
    + registers an object type with the names of its exposed fields
    # schema
  graphql.register_resolver
    fn (state: schema_state, type_name: string, field_name: string, resolver_id: string) -> result[schema_state, string]
    + associates a resolver id with a (type, field) pair
    - returns error when the type or field is not registered
    # schema
  graphql.parse_query
    fn (query: string) -> result[query_ast, string]
    + parses a GraphQL query into a selection-set AST
    - returns error on syntax errors
    # parsing
  graphql.execute
    fn (state: schema_state, ast: query_ast, root_type: string) -> result[string, string]
    + walks the selection set, invoking registered resolvers and returning a JSON response
    - returns error when a selected field has no resolver
    # execution
    -> std.json.encode_object
  graphql.format_error
    fn (message: string, path: list[string]) -> string
    + formats a GraphQL error with a message and field path
    # errors
