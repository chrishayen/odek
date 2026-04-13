# Requirement: "a library that generates database access code from a SQL schema file or a database connection string"

Parses SQL DDL (or introspects a live database), builds a table model, and emits code for multiple target database access layers.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full text of a file
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + creates or overwrites the file
      - returns error when the parent directory does not exist
      # filesystem
  std.db
    std.db.connect
      @ (dsn: string) -> result[db_handle, string]
      + opens a connection using a driver inferred from the dsn scheme
      - returns error on malformed dsn
      # database
    std.db.list_tables
      @ (db: db_handle) -> result[list[string], string]
      + returns the names of user tables in the current schema
      # database
    std.db.describe_table
      @ (db: db_handle, table: string) -> result[list[column_info], string]
      + returns column name, sql type, nullability, and primary-key flag for each column
      - returns error when the table is not found
      # database

sqlgen
  sqlgen.parse_schema_sql
    @ (sql: string) -> result[list[table_model], string]
    + parses CREATE TABLE statements into structured table models
    + records column name, type, nullability, and primary key
    - returns error on unterminated statements or unknown keywords
    # sql_parsing
  sqlgen.schema_from_dsn
    @ (dsn: string) -> result[list[table_model], string]
    + connects and builds table models by introspection
    - returns error when the connection fails
    # introspection
    -> std.db.connect
    -> std.db.list_tables
    -> std.db.describe_table
  sqlgen.load_schema
    @ (source: string) -> result[list[table_model], string]
    + dispatches to file parsing when source looks like a path and to dsn introspection otherwise
    + reads file contents before parsing
    # dispatch
    -> std.fs.read_all
  sqlgen.map_sql_type
    @ (sql_type: string, nullable: bool) -> string
    + returns the language-agnostic field type for a given sql column type
    + wraps the type in an optional marker when nullable
    # type_mapping
  sqlgen.render_orm_a
    @ (tables: list[table_model]) -> string
    + emits code for orm style "a" with field tags
    # codegen
  sqlgen.render_orm_b
    @ (tables: list[table_model]) -> string
    + emits code for orm style "b" with struct tags and a table registration block
    # codegen
  sqlgen.render_query_builder
    @ (tables: list[table_model]) -> string
    + emits code for a query-builder style layer with typed column constants
    # codegen
  sqlgen.render_raw_sql_helpers
    @ (tables: list[table_model]) -> string
    + emits raw-sql helper functions (select, insert, update, delete) per table
    # codegen
  sqlgen.generate
    @ (source: string, target: string, out_dir: string) -> result[void, string]
    + loads schema, renders the requested target, writes the result to disk
    - returns error when target is not one of the supported names
    # orchestration
    -> std.fs.write_all
