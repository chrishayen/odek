# Requirement: "a type-safe code generator from SQL queries"

Parses SQL schema and query files, resolves query parameters and result columns against the schema, then emits typed function signatures.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file's contents as a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_all
      fn (path: string, content: string) -> result[void, string]
      + writes the content, creating or truncating the file
      - returns error when the path is not writable
      # filesystem
  std.strings
    std.strings.split_statements
      fn (sql: string) -> list[string]
      + splits a SQL text into top-level statements on semicolons outside quoted regions
      # text_processing
    std.strings.to_pascal_case
      fn (snake: string) -> string
      + converts snake_case identifiers to PascalCase
      # text_processing

sqlgen
  sqlgen.parse_schema
    fn (sql: string) -> result[schema, string]
    + parses CREATE TABLE statements into a schema with tables and typed columns
    - returns error on unknown column type
    # schema_parsing
    -> std.strings.split_statements
  sqlgen.parse_query
    fn (sql: string) -> result[query_ast, string]
    + extracts the query kind, referenced tables, placeholders, and selected columns
    - returns error when the query references an unknown clause form
    # query_parsing
  sqlgen.resolve_query
    fn (q: query_ast, s: schema) -> result[typed_query, string]
    + binds each placeholder to a column type and resolves each result column to its table type
    - returns error when a referenced table or column does not exist in the schema
    # type_resolution
  sqlgen.map_sql_type
    fn (sql_type: string) -> result[string, string]
    + maps a SQL column type name to a language-neutral type name like i32, string, bool
    - returns error for unsupported SQL types
    # type_mapping
  sqlgen.emit_signature
    fn (name: string, tq: typed_query) -> string
    + renders a typed function signature and parameter list for the query
    # code_emission
    -> std.strings.to_pascal_case
    -> sqlgen.map_sql_type
  sqlgen.generate
    fn (schema_sql: string, queries: map[string, string]) -> result[string, string]
    + parses schema and each named query, returns concatenated typed signatures
    - returns error if any query fails to resolve against the schema
    # pipeline
    -> sqlgen.parse_schema
    -> sqlgen.parse_query
    -> sqlgen.resolve_query
    -> sqlgen.emit_signature
