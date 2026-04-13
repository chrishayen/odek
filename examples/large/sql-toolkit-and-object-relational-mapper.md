# Requirement: "a SQL toolkit and object-relational mapper"

Provides a typed schema definition, a query builder, a session that tracks dirty objects, and a pluggable SQL driver.

std
  std.db
    std.db.open
      @ (url: string) -> result[db_conn, string]
      + opens a database connection from a URL
      - returns error when the URL is invalid or the server unreachable
      # database
    std.db.execute
      @ (conn: db_conn, sql: string, params: list[sql_value]) -> result[i64, string]
      + runs a statement and returns the number of rows affected
      - returns error on SQL or constraint failure
      # database
    std.db.query
      @ (conn: db_conn, sql: string, params: list[sql_value]) -> result[list[row], string]
      + runs a query and returns the result rows
      - returns error on SQL failure
      # database

orm
  orm.define_table
    @ (name: string, columns: list[column_spec]) -> table_def
    + returns a table definition with the given columns
    # schema
  orm.primary_key
    @ (columns: list[string]) -> constraint
    + returns a primary-key constraint over the named columns
    # schema
  orm.foreign_key
    @ (columns: list[string], ref_table: string, ref_columns: list[string]) -> constraint
    + returns a foreign-key constraint referencing another table
    # schema
  orm.create_schema
    @ (conn: db_conn, tables: list[table_def]) -> result[void, string]
    + emits DDL to create every table and its constraints
    - returns error when a referenced table is missing
    # schema
    -> std.db.execute
  orm.select
    @ (table: string) -> query
    + returns a new select query targeting the given table
    # query_builder
  orm.where_eq
    @ (q: query, column: string, value: sql_value) -> query
    + returns a query with an additional equality predicate
    # query_builder
  orm.join
    @ (q: query, table: string, on_left: string, on_right: string) -> query
    + returns a query joined to another table on the given columns
    # query_builder
  orm.order_by
    @ (q: query, column: string, direction: order_dir) -> query
    + returns a query ordered by the given column and direction
    # query_builder
  orm.to_sql
    @ (q: query) -> tuple[string, list[sql_value]]
    + returns the parameterized SQL and bind values for the query
    # query_builder
  orm.new_session
    @ (conn: db_conn) -> session_state
    + creates a session bound to a connection and an empty identity map
    # session
  orm.session_add
    @ (state: session_state, table: string, entity: map[string, sql_value]) -> session_state
    + stages a new row for insertion on the next flush
    # session
  orm.session_get
    @ (state: session_state, table: string, key: map[string, sql_value]) -> result[optional[map[string, sql_value]], string]
    + returns a row from the identity map, loading it from the database on miss
    - returns error on database failure
    # session
    -> std.db.query
  orm.session_flush
    @ (state: session_state) -> result[session_state, string]
    + emits pending inserts, updates, and deletes in dependency order
    - returns error on constraint violation
    # session
    -> std.db.execute
  orm.begin_transaction
    @ (state: session_state) -> result[session_state, string]
    + opens a database transaction for subsequent flushes
    - returns error when a transaction is already active
    # transaction
  orm.commit
    @ (state: session_state) -> result[session_state, string]
    + commits the active transaction
    - returns error when no transaction is active
    # transaction
  orm.rollback
    @ (state: session_state) -> result[session_state, string]
    + rolls back the active transaction and clears pending changes
    # transaction
