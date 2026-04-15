# Requirement: "a GraphQL server and client library with schema-driven code generation and REST endpoints"

The project layer builds a schema, resolves queries, serves an HTTP endpoint, and offers a parallel REST facade. Code generation is a separate helper that emits bindings from a schema.

std
  std.http
    std.http.listen
      fn (addr: string, handler: request_handler) -> result[void, string]
      + binds to addr and dispatches requests to handler
      - returns error when the address is unavailable
      # http
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a flat JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

gqlkit
  gqlkit.new_schema
    fn () -> gql_schema
    + creates an empty schema
    # schema
  gqlkit.add_type
    fn (schema: gql_schema, name: string, fields: map[string,string]) -> gql_schema
    + registers a type with named fields and scalar field types
    # schema
  gqlkit.add_resolver
    fn (schema: gql_schema, type_name: string, field: string, resolver: resolver_handle) -> gql_schema
    + attaches a resolver to a type/field pair
    # resolvers
  gqlkit.parse_query
    fn (source: string) -> result[gql_query, string]
    + parses a GraphQL query into an operation tree
    - returns error on invalid syntax
    # parsing
  gqlkit.execute
    fn (schema: gql_schema, query: gql_query, variables: map[string,string]) -> result[string, string]
    + runs resolvers and returns a JSON response
    - returns error when a referenced field has no resolver
    # execution
    -> std.json.encode_object
  gqlkit.serve_graphql
    fn (schema: gql_schema, addr: string) -> result[void, string]
    + listens on addr and handles POST /graphql requests
    # serving
    -> std.http.listen
    -> std.json.parse_object
  gqlkit.serve_rest
    fn (schema: gql_schema, addr: string) -> result[void, string]
    + exposes schema types as REST resources on addr
    # serving
    -> std.http.listen
  gqlkit.generate_bindings
    fn (schema: gql_schema) -> string
    + emits client-side binding source from the schema
    ? output is a generic pseudo-code form; target language is caller's problem
    # codegen
