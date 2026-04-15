# Requirement: "a database migration generator that derives SQL migrations by diffing a model against an existing schema"

Model and existing SQL are parsed into a common schema representation; the diff yields up/down SQL.

std
  std.io
    std.io.read_file
      fn (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the path does not exist
      # filesystem

migrate_gen
  migrate_gen.parse_sql_schema
    fn (sql: string) -> result[schema, string]
    + parses CREATE TABLE statements into a schema of tables and columns
    - returns error on unterminated statements
    # parsing
  migrate_gen.parse_model
    fn (model_source: string) -> result[schema, string]
    + parses a struct-style model description into a schema
    - returns error when a field has no type annotation
    # parsing
  migrate_gen.diff
    fn (from: schema, to: schema) -> schema_diff
    + returns a diff listing added, dropped, and altered tables and columns
    + empty diff when schemas are identical
    # diffing
  migrate_gen.render_up
    fn (diff: schema_diff) -> string
    + renders the forward migration as SQL statements
    # codegen
  migrate_gen.render_down
    fn (diff: schema_diff) -> string
    + renders the reverse migration as SQL statements
    # codegen
