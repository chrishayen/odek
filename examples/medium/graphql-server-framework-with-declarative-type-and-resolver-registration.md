# Requirement: "a GraphQL server framework with declarative type and resolver registration"

Users register object types and resolvers, then execute GraphQL queries against them.

std: (all units exist)

gql
  gql.new_schema
    fn () -> schema_state
    + creates an empty schema with no types and no root resolvers
    # construction
  gql.object
    fn (s: schema_state, name: string, fields: list[field_def]) -> schema_state
    + registers an object type with its fields
    ? each field_def has a name, a GraphQL type, and a resolver fn(parent, args) -> result[value, string]
    # type_registration
  gql.enum
    fn (s: schema_state, name: string, values: list[string]) -> schema_state
    + registers an enum type
    # type_registration
  gql.query_root
    fn (s: schema_state, type_name: string) -> schema_state
    + marks type_name as the query root object
    # root
  gql.mutation_root
    fn (s: schema_state, type_name: string) -> schema_state
    + marks type_name as the mutation root object
    # root
  gql.parse_query
    fn (source: string) -> result[query_doc, string]
    + parses a GraphQL query string into a document
    - returns error on syntax issues
    # parsing
  gql.validate_query
    fn (s: schema_state, doc: query_doc) -> result[void, string]
    + verifies that every selected field exists on its parent type and that arguments match
    - returns error describing the first mismatch
    # validation
  gql.execute
    fn (s: schema_state, doc: query_doc, variables: map[string, value]) -> result[value, string]
    + runs resolvers over the selection set and returns the resulting tree
    - returns error on validation or resolver failure
    # execution
  gql.run
    fn (s: schema_state, source: string, variables: map[string, value]) -> result[value, string]
    + one-shot parse-validate-execute entry point
    # top_level
