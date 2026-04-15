# Requirement: "a query builder that produces compile-time verified queries from generated schema metadata"

Schema-driven query construction. The project layer models tables, columns, and type-checked expressions then renders SQL.

std: (all units exist)

typed_query
  typed_query.define_table
    fn (name: string, columns: list[column_def]) -> table_schema
    + creates a table schema from name and typed column definitions
    # schema
  typed_query.column_ref
    fn (table: table_schema, column_name: string) -> result[column_ref, string]
    + returns a typed reference to the column
    - returns error when the column does not exist on the table
    # schema_lookup
  typed_query.eq
    fn (left: column_ref, right: value) -> result[predicate, string]
    + builds an equality predicate
    - returns error when value type does not match the column type
    # predicates
  typed_query.and
    fn (a: predicate, b: predicate) -> predicate
    + combines two predicates with logical AND
    # predicates
  typed_query.select
    fn (table: table_schema, columns: list[column_ref]) -> query_builder
    + starts a SELECT over the given table and columns
    + validates that each column belongs to the table
    # query_construction
  typed_query.where
    fn (q: query_builder, pred: predicate) -> query_builder
    + attaches a WHERE clause
    # query_construction
  typed_query.render
    fn (q: query_builder) -> tuple[string, list[value]]
    + returns parameterized SQL and the ordered bound values
    ? placeholders are rendered as $1, $2, ...
    # rendering
  typed_query.insert
    fn (table: table_schema, values: map[string,value]) -> result[tuple[string, list[value]], string]
    + builds an INSERT statement and bindings
    - returns error when a required column is missing or has the wrong type
    # query_construction
