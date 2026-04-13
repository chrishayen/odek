# Requirement: "discover sql schemas and convert them to avro schemas, and query sql records into avro bytes"

Introspects a sql database, converts table definitions to avro record schemas, and streams rows as avro-encoded bytes.

std
  std.sql
    std.sql.connect
      @ (dsn: string) -> result[sql_conn, string]
      + opens a connection described by dsn
      - returns error when connection fails
      # database
    std.sql.query
      @ (conn: sql_conn, sql: string) -> result[row_cursor, string]
      + executes a query returning a cursor over rows
      # database
    std.sql.next_row
      @ (cursor: row_cursor) -> optional[list[sql_value]]
      + returns the next row or none when exhausted
      # database
    std.sql.list_tables
      @ (conn: sql_conn) -> result[list[string], string]
      + returns all user-visible table names
      # database
    std.sql.describe_table
      @ (conn: sql_conn, table: string) -> result[list[sql_column], string]
      + returns column names and sql type info for the given table
      - returns error when table does not exist
      # database
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + encodes a json value as a string
      # serialization

sql_to_avro
  sql_to_avro.discover_schemas
    @ (conn: sql_conn) -> result[list[avro_schema], string]
    + returns an avro record schema for every table in the database
    # discovery
    -> std.sql.list_tables
    -> std.sql.describe_table
  sql_to_avro.table_to_avro
    @ (table: string, columns: list[sql_column]) -> avro_schema
    + builds an avro record schema with one field per column
    + maps sql types to avro primitives; nullable columns become unions with null
    ? decimal and date types are encoded as logical types
    # schema_translation
  sql_to_avro.avro_schema_to_json
    @ (schema: avro_schema) -> string
    + serializes an avro schema to its canonical json form
    # schema_translation
    -> std.json.encode
  sql_to_avro.query_to_avro_bytes
    @ (conn: sql_conn, sql: string, schema: avro_schema) -> result[bytes, string]
    + executes sql and returns the rows as an avro object container file
    - returns error when a row does not match the given schema
    # encoding
    -> std.sql.query
    -> std.sql.next_row
  sql_to_avro.encode_row
    @ (row: list[sql_value], schema: avro_schema) -> result[bytes, string]
    + encodes a single row using avro binary encoding
    - returns error when a value cannot be coerced to its declared type
    # encoding
