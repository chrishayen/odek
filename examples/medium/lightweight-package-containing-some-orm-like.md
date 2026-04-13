# Requirement: "an ORM-lite helper for an embedded SQL database"

Provides struct-to-row mapping, basic CRUD, and simple where clauses against a single-file SQL database.

std
  std.sql
    std.sql.open
      @ (path: string) -> result[sql_conn, string]
      + opens or creates a single-file database
      - returns error when the path is not writable
      # database
    std.sql.exec
      @ (conn: sql_conn, statement: string, params: list[string]) -> result[i64, string]
      + executes a parameterized statement and returns rows affected
      - returns error on invalid SQL
      # database
    std.sql.query_rows
      @ (conn: sql_conn, statement: string, params: list[string]) -> result[list[map[string, string]], string]
      + returns rows as maps of column-to-text-value
      - returns error on invalid SQL
      # database

orm_lite
  orm_lite.define_table
    @ (name: string, columns: list[column_def], primary_key: string) -> table_def
    + captures a table definition used by subsequent operations
    # schema
  orm_lite.create_table
    @ (conn: sql_conn, t: table_def) -> result[void, string]
    + issues CREATE TABLE IF NOT EXISTS for the definition
    - returns error on invalid column type
    # schema
    -> std.sql.exec
  orm_lite.insert
    @ (conn: sql_conn, t: table_def, row: map[string, string]) -> result[i64, string]
    + inserts a row and returns the last row id
    - returns error on constraint violation
    # write
    -> std.sql.exec
  orm_lite.find_by
    @ (conn: sql_conn, t: table_def, column: string, value: string) -> result[list[map[string, string]], string]
    + returns all rows where the column matches the value
    # read
    -> std.sql.query_rows
  orm_lite.update_by_pk
    @ (conn: sql_conn, t: table_def, row: map[string, string]) -> result[i64, string]
    + updates the row identified by primary key and returns rows affected
    - returns error when the row lacks a primary key value
    # write
    -> std.sql.exec
  orm_lite.delete_by_pk
    @ (conn: sql_conn, t: table_def, pk_value: string) -> result[i64, string]
    + deletes the row matching the primary key and returns rows affected
    # write
    -> std.sql.exec
