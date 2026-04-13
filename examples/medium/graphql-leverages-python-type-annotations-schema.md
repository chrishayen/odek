# Requirement: "a GraphQL schema library driven by typed record declarations"

Turn a set of typed record declarations into a GraphQL schema and execute queries against it.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a tagged value
      - returns error on malformed input
      # serialization
    std.json.encode_value
      @ (value: json_value) -> string
      + encodes a tagged value as JSON
      # serialization

typed_graphql
  typed_graphql.new_schema
    @ () -> schema_state
    + creates an empty schema
    # construction
  typed_graphql.register_type
    @ (schema: schema_state, name: string, fields: map[string, string]) -> result[schema_state, string]
    + registers an object type with the given field_name -> type_name map
    - returns error when the name collides with a built-in scalar
    - returns error when a referenced type is undefined after all registrations
    # schema
  typed_graphql.register_query
    @ (schema: schema_state, name: string, return_type: string, resolve: closure[map[string, string]]) -> schema_state
    + exposes a top-level query field with a resolver
    # schema
  typed_graphql.parse_query
    @ (source: string) -> result[query_ast, string]
    + parses a GraphQL query document into an AST
    - returns error on syntax failure
    # parsing
  typed_graphql.validate
    @ (schema: schema_state, query: query_ast) -> result[void, list[string]]
    + checks selection sets, field names, and argument types against the schema
    - returns the collected error messages when validation fails
    # validation
  typed_graphql.execute
    @ (schema: schema_state, query: query_ast, variables: map[string, string]) -> result[string, string]
    + runs the query and returns the JSON-encoded result document
    - returns error when a resolver raises
    # execution
    -> std.json.parse_value
    -> std.json.encode_value
