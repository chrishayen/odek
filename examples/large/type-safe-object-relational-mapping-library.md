# Requirement: "a type-safe object-relational mapping library"

Maps records to database rows. Users declare a schema, then insert, find, update, and delete through the ORM.

std
  std.sql
    std.sql.execute
      fn (conn: db_conn, sql: string, params: list[string]) -> result[i64, string]
      + executes a statement and returns affected rows
      - returns error on sql error
      # sql
    std.sql.query
      fn (conn: db_conn, sql: string, params: list[string]) -> result[list[map[string, string]], string]
      + returns rows as string-to-string maps
      - returns error on sql error
      # sql
  std.strings
    std.strings.join
      fn (parts: list[string], sep: string) -> string
      + joins parts with sep
      # strings

orm
  orm.new_schema
    fn (table: string) -> schema_state
    + returns an empty schema bound to the given table
    # construction
  orm.add_column
    fn (schema: schema_state, name: string, sql_type: string) -> schema_state
    + appends a column definition
    + preserves insertion order
    # schema
  orm.mark_primary_key
    fn (schema: schema_state, name: string) -> result[schema_state, string]
    + designates the column as the primary key
    - returns error when the column is not in the schema
    # schema
  orm.create_table_sql
    fn (schema: schema_state) -> string
    + returns a CREATE TABLE statement for the schema
    # codegen
    -> std.strings.join
  orm.insert
    fn (conn: db_conn, schema: schema_state, record: map[string, string]) -> result[i64, string]
    + inserts the record and returns the number of affected rows
    - returns error when required columns are missing
    - returns error on sql failure
    # write
    -> std.sql.execute
  orm.find_by_id
    fn (conn: db_conn, schema: schema_state, id: string) -> result[optional[map[string, string]], string]
    + returns the row with the matching primary key
    - returns none when no row matches
    - returns error on sql failure
    # read
    -> std.sql.query
  orm.find_where
    fn (conn: db_conn, schema: schema_state, column: string, value: string) -> result[list[map[string, string]], string]
    + returns rows whose column equals the given value
    - returns error on sql failure
    - returns error when the column is not in the schema
    # read
    -> std.sql.query
  orm.update
    fn (conn: db_conn, schema: schema_state, id: string, changes: map[string, string]) -> result[i64, string]
    + updates the primary-keyed row and returns affected rows
    - returns error when changes contains unknown columns
    - returns error on sql failure
    # write
    -> std.sql.execute
  orm.delete_by_id
    fn (conn: db_conn, schema: schema_state, id: string) -> result[i64, string]
    + deletes the primary-keyed row and returns affected rows
    - returns error on sql failure
    # write
    -> std.sql.execute
