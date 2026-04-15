# Requirement: "an ORM layered on a SQL query builder"

Maps row records to and from models and exposes CRUD operations built on a lower-level query builder.

std
  std.sql
    std.sql.exec
      fn (conn: sql_conn, statement: string, params: list[string]) -> result[i64, string]
      + executes a parameterized statement and returns rows affected
      - returns error on invalid SQL
      # database
    std.sql.query_rows
      fn (conn: sql_conn, statement: string, params: list[string]) -> result[list[map[string, string]], string]
      + returns rows as maps of column-to-text-value
      - returns error on invalid SQL
      # database

orm
  orm.define_model
    fn (table: string, columns: list[string], primary_key: string) -> model
    + returns a model description used by subsequent operations
    # schema
  orm.build_select
    fn (m: model, where_columns: list[string]) -> string
    + returns a SELECT statement with placeholders for the where columns
    # query_building
  orm.build_insert
    fn (m: model, columns: list[string]) -> string
    + returns an INSERT statement with placeholders
    # query_building
  orm.build_update
    fn (m: model, set_columns: list[string]) -> string
    + returns an UPDATE ... WHERE <pk> = ? statement
    # query_building
  orm.find
    fn (conn: sql_conn, m: model, pk_value: string) -> result[optional[map[string, string]], string]
    + returns the row matching the primary key, or none when absent
    - returns error on query failure
    # read
    -> std.sql.query_rows
  orm.save
    fn (conn: sql_conn, m: model, row: map[string, string]) -> result[void, string]
    + inserts the row, or updates it when the primary key is already present
    - returns error on constraint violation
    # write
    -> std.sql.exec
    -> std.sql.query_rows
  orm.delete
    fn (conn: sql_conn, m: model, pk_value: string) -> result[i64, string]
    + returns the number of rows removed
    - returns error on query failure
    # write
    -> std.sql.exec
