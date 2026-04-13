# Requirement: "a database migration library"

Applies and rolls back versioned SQL migrations against a database connection, tracking applied versions in a schema table.

std
  std.fs
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns filenames in the directory sorted lexicographically
      - returns error when the directory does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full contents of a text file
      - returns error when the file is missing
      # filesystem
  std.sql
    std.sql.exec
      @ (conn: sql_conn, statement: string) -> result[i64, string]
      + executes a statement and returns rows affected
      - returns error on invalid SQL or constraint violation
      # database
    std.sql.query_strings
      @ (conn: sql_conn, statement: string) -> result[list[string], string]
      + returns the first column of every row
      - returns error on invalid SQL
      # database

migrations
  migrations.load_from_dir
    @ (dir: string) -> result[list[migration], string]
    + loads migration files named "<version>_<name>.<up|down>.sql" from a directory
    - returns error when version prefixes are not sortable integers
    # loading
    -> std.fs.list_dir
    -> std.fs.read_all
  migrations.ensure_schema
    @ (conn: sql_conn) -> result[void, string]
    + creates the tracking table if it does not exist
    # bootstrap
    -> std.sql.exec
  migrations.applied_versions
    @ (conn: sql_conn) -> result[list[string], string]
    + returns versions already applied, in ascending order
    # state
    -> std.sql.query_strings
  migrations.up
    @ (conn: sql_conn, all_migrations: list[migration]) -> result[i32, string]
    + applies every pending migration in order and returns the count applied
    - returns error when an up statement fails
    # apply
    -> std.sql.exec
  migrations.down
    @ (conn: sql_conn, all_migrations: list[migration], steps: i32) -> result[i32, string]
    + rolls back the most recent N migrations
    - returns error when a down statement fails
    # rollback
    -> std.sql.exec
