# Requirement: "a GraphQL query execution engine"

Accept a schema, parse and validate incoming queries, and execute them against registered resolvers. Transport is the caller's concern.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

graphql_engine
  graphql_engine.new_schema
    @ () -> schema_state
    + creates an empty schema
    # construction
  graphql_engine.add_field
    @ (schema: schema_state, type_name: string, field_name: string, return_type: string, resolve: closure[string]) -> schema_state
    + registers a field with its resolver on a type
    # schema
  graphql_engine.parse
    @ (source: string) -> result[query_ast, string]
    + parses a query document
    - returns error on syntax failure
    # parsing
  graphql_engine.validate
    @ (schema: schema_state, ast: query_ast) -> result[void, list[string]]
    + checks every selected field exists on its type
    - returns the list of validation errors
    # validation
  graphql_engine.execute
    @ (schema: schema_state, ast: query_ast, variables: map[string, string]) -> result[string, string]
    + executes the query and returns a JSON response body
    - returns error when a resolver raises
    # execution
    -> std.json.parse_object
    -> std.json.encode_object
