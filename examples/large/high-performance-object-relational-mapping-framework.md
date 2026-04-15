# Requirement: "a high-performance object-relational mapping framework"

Maps records between rows and struct-like maps, builds parameterized SQL for common operations, and executes queries through a pluggable connection.

std
  std.sql
    std.sql.open
      fn (dsn: string) -> result[sql_conn, string]
      + opens a connection to a database using a data source name
      - returns error on invalid dsn
      # database
    std.sql.exec
      fn (conn: sql_conn, query: string, args: list[bytes]) -> result[i64, string]
      + executes a statement and returns the affected row count
      - returns error on syntax or constraint failure
      # database
    std.sql.query
      fn (conn: sql_conn, query: string, args: list[bytes]) -> result[list[map[string, bytes]], string]
      + executes a query and returns a list of row maps
      - returns error when the driver reports a failure
      # database
    std.sql.close
      fn (conn: sql_conn) -> result[void, string]
      + releases resources held by the connection
      # database
  std.strings
    std.strings.join
      fn (parts: list[string], sep: string) -> string
      + joins parts with the separator
      # strings

orm
  orm.define_entity
    fn (table: string, fields: list[tuple[string, string]]) -> entity_schema
    + defines a schema binding a table name to ordered (column, type) pairs
    # schema
  orm.build_select
    fn (schema: entity_schema, where: map[string, bytes]) -> tuple[string, list[bytes]]
    + returns a parameterized SELECT and the bound argument list
    + returns a SELECT with no WHERE clause when the map is empty
    # query_building
    -> std.strings.join
  orm.build_insert
    fn (schema: entity_schema, row: map[string, bytes]) -> result[tuple[string, list[bytes]], string]
    + returns a parameterized INSERT and the bound argument list
    - returns error when row references a column not in the schema
    # query_building
    -> std.strings.join
  orm.build_update
    fn (schema: entity_schema, set: map[string, bytes], where: map[string, bytes]) -> result[tuple[string, list[bytes]], string]
    + returns a parameterized UPDATE and the bound argument list
    - returns error when set is empty
    # query_building
    -> std.strings.join
  orm.build_delete
    fn (schema: entity_schema, where: map[string, bytes]) -> tuple[string, list[bytes]]
    + returns a parameterized DELETE and the bound argument list
    # query_building
  orm.find_one
    fn (conn: sql_conn, schema: entity_schema, where: map[string, bytes]) -> result[optional[map[string, bytes]], string]
    + returns the first matching row as a column-value map
    - returns empty when no row matches
    - returns error when the query fails
    # execution
    -> std.sql.query
  orm.find_all
    fn (conn: sql_conn, schema: entity_schema, where: map[string, bytes]) -> result[list[map[string, bytes]], string]
    + returns all matching rows
    - returns error when the query fails
    # execution
    -> std.sql.query
  orm.save
    fn (conn: sql_conn, schema: entity_schema, row: map[string, bytes]) -> result[i64, string]
    + inserts the row and returns affected count
    - returns error when the row is missing a required column
    # execution
    -> std.sql.exec
  orm.update
    fn (conn: sql_conn, schema: entity_schema, set: map[string, bytes], where: map[string, bytes]) -> result[i64, string]
    + applies updates to matching rows and returns affected count
    - returns error when set is empty
    # execution
    -> std.sql.exec
  orm.remove
    fn (conn: sql_conn, schema: entity_schema, where: map[string, bytes]) -> result[i64, string]
    + deletes matching rows and returns affected count
    # execution
    -> std.sql.exec
