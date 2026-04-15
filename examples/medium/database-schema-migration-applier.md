# Requirement: "a library for embedding and applying ordered database schema migrations"

Migrations are registered as data; the library tracks applied versions and runs pending ones in order inside a transaction.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits on every separator occurrence
      # strings
  std.sql
    std.sql.exec
      fn (conn: sql_conn, query: string, args: list[string]) -> result[i64, string]
      + executes a statement and returns rows affected
      - returns error on invalid SQL
      # database
    std.sql.query
      fn (conn: sql_conn, query: string, args: list[string]) -> result[list[list[string]], string]
      + executes a query and returns rows as string lists
      # database
    std.sql.begin_tx
      fn (conn: sql_conn) -> result[sql_tx, string]
      + opens a transaction
      # database
    std.sql.commit
      fn (tx: sql_tx) -> result[void, string]
      + commits the transaction
      # database
    std.sql.rollback
      fn (tx: sql_tx) -> result[void, string]
      + rolls the transaction back
      # database

migrations
  migrations.new
    fn () -> migration_set
    + creates an empty migration set
    # construction
  migrations.add
    fn (set: migration_set, version: i64, name: string, up_sql: string, down_sql: string) -> migration_set
    + appends a migration to the set
    - returns an error-marked set when version is not strictly greater than the previous entry
    # schema
  migrations.ensure_table
    fn (conn: sql_conn) -> result[void, string]
    + creates the schema_migrations tracking table if it does not exist
    # database
    -> std.sql.exec
  migrations.applied_versions
    fn (conn: sql_conn) -> result[list[i64], string]
    + returns versions currently recorded as applied, sorted ascending
    # query
    -> std.sql.query
  migrations.up
    fn (conn: sql_conn, set: migration_set) -> result[i32, string]
    + applies all pending migrations within a transaction and returns the count applied
    - rolls back and returns error when any migration fails
    # migration
    -> std.sql.begin_tx
    -> std.sql.exec
    -> std.sql.commit
    -> std.sql.rollback
  migrations.down
    fn (conn: sql_conn, set: migration_set, target_version: i64) -> result[i32, string]
    + reverts applied migrations with version greater than target_version in reverse order
    - rolls back on failure mid-run
    # migration
    -> std.sql.begin_tx
    -> std.sql.exec
    -> std.sql.commit
    -> std.sql.rollback
  migrations.status
    fn (conn: sql_conn, set: migration_set) -> result[list[migration_status_entry], string]
    + returns every migration annotated as applied or pending
    # query
