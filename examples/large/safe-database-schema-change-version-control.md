# Requirement: "a database schema migration library with versioning, validation, and safe apply/rollback"

Migrations are ordered records with up and down scripts. The core tracks applied state in a meta table and enforces a strict forward/backward order.

std
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + computes SHA-256 and returns a lowercase hex digest
      # hashing
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.db
    std.db.exec
      @ (conn: db_conn, sql: string) -> result[void, string]
      + executes a statement that returns no rows
      - returns error on SQL failure
      # database
    std.db.query_rows
      @ (conn: db_conn, sql: string) -> result[list[map[string,string]], string]
      + executes a query and returns rows as string maps
      - returns error on SQL failure
      # database
    std.db.begin_tx
      @ (conn: db_conn) -> result[db_tx, string]
      + opens a transaction
      # database
    std.db.commit_tx
      @ (tx: db_tx) -> result[void, string]
      + commits a transaction
      # database
    std.db.rollback_tx
      @ (tx: db_tx) -> result[void, string]
      + rolls back a transaction
      # database

migrate
  migrate.migration_load
    @ (version: i64, name: string, up_sql: string, down_sql: string) -> migration
    + constructs a migration record stamped with a content checksum
    # construction
    -> std.hash.sha256_hex
  migrate.plan_sort
    @ (migrations: list[migration]) -> result[list[migration], string]
    + returns migrations sorted ascending by version
    - returns error when two migrations share a version
    # planning
  migrate.ensure_meta_table
    @ (conn: db_conn) -> result[void, string]
    + creates the schema_migrations table if it does not exist
    # bootstrap
    -> std.db.exec
  migrate.applied_versions
    @ (conn: db_conn) -> result[list[i64], string]
    + returns versions already recorded in the meta table
    # state
    -> std.db.query_rows
  migrate.pending
    @ (sorted: list[migration], applied: list[i64]) -> list[migration]
    + returns migrations not yet applied
    # planning
  migrate.verify_checksums
    @ (sorted: list[migration], conn: db_conn) -> result[void, string]
    + verifies that applied migrations still match their stored checksums
    - returns error when a previously applied migration's content has changed
    # integrity
    -> std.db.query_rows
  migrate.apply_one
    @ (conn: db_conn, m: migration) -> result[void, string]
    + runs the up script inside a transaction and records the version and checksum
    - rolls back and returns error on any SQL failure
    # apply
    -> std.db.begin_tx
    -> std.db.exec
    -> std.db.commit_tx
    -> std.db.rollback_tx
    -> std.time.now_millis
  migrate.apply_all
    @ (conn: db_conn, migrations: list[migration]) -> result[i32, string]
    + applies all pending migrations in order and returns the number applied
    - returns error on checksum mismatch or SQL failure
    # apply
  migrate.rollback_one
    @ (conn: db_conn, m: migration) -> result[void, string]
    + runs the down script inside a transaction and removes the version row
    - returns error when the migration has not been applied
    # rollback
    -> std.db.begin_tx
    -> std.db.exec
    -> std.db.commit_tx
    -> std.db.rollback_tx
  migrate.rollback_to
    @ (conn: db_conn, migrations: list[migration], target_version: i64) -> result[i32, string]
    + rolls back migrations in reverse order until the applied head reaches target_version
    - returns error on SQL failure or when target_version is not in the history
    # rollback
