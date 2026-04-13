# Requirement: "a key-value store with redis-style commands backed by a relational database"

Exposes string, list, hash, and set operations, persisting state through a sql-backed storage layer. Project layer owns the command semantics; std provides sql and time.

std
  std.sql
    std.sql.open
      @ (dsn: string) -> result[db_handle, string]
      + opens a connection pool to the database at dsn
      - returns error when the database cannot be reached
      # database
    std.sql.exec
      @ (db: db_handle, query: string, args: list[sql_value]) -> result[i64, string]
      + executes a statement and returns rows affected
      - returns error on syntax or constraint failure
      # database
    std.sql.query_rows
      @ (db: db_handle, query: string, args: list[sql_value]) -> result[list[sql_row], string]
      + runs a query and returns all rows
      # database
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

kvstore
  kvstore.open
    @ (db_path: string) -> result[store_state, string]
    + opens or creates the backing database and initializes schema
    - returns error when the file cannot be created
    # construction
    -> std.sql.open
    -> std.sql.exec
  kvstore.set_string
    @ (s: store_state, key: string, value: string, ttl_ms: i64) -> result[void, string]
    + stores a string value for key with an optional ttl (0 for no expiry)
    # strings
    -> std.sql.exec
    -> std.time.now_millis
  kvstore.get_string
    @ (s: store_state, key: string) -> result[optional[string], string]
    + returns the string value or none when absent or expired
    # strings
    -> std.sql.query_rows
    -> std.time.now_millis
  kvstore.delete
    @ (s: store_state, key: string) -> result[i64, string]
    + removes the key and returns the number of keys deleted
    # keyspace
    -> std.sql.exec
  kvstore.exists
    @ (s: store_state, key: string) -> result[bool, string]
    + returns true when the key exists and is not expired
    # keyspace
  kvstore.list_push
    @ (s: store_state, key: string, value: string, left: bool) -> result[i64, string]
    + pushes value onto the left or right of a list and returns the new length
    # lists
    -> std.sql.exec
  kvstore.list_pop
    @ (s: store_state, key: string, left: bool) -> result[optional[string], string]
    + pops and returns an element from the list, or none when empty
    # lists
    -> std.sql.exec
  kvstore.hash_set
    @ (s: store_state, key: string, field: string, value: string) -> result[void, string]
    + sets a field in the hash
    # hashes
    -> std.sql.exec
  kvstore.hash_get
    @ (s: store_state, key: string, field: string) -> result[optional[string], string]
    + returns the field value or none when absent
    # hashes
    -> std.sql.query_rows
  kvstore.set_add
    @ (s: store_state, key: string, member: string) -> result[bool, string]
    + adds a member to the set and returns true when newly added
    # sets
    -> std.sql.exec
  kvstore.set_members
    @ (s: store_state, key: string) -> result[list[string], string]
    + returns all members of the set
    # sets
    -> std.sql.query_rows
  kvstore.expire_sweep
    @ (s: store_state) -> result[i64, string]
    + removes all expired keys and returns the count
    # ttl
    -> std.sql.exec
    -> std.time.now_millis
  kvstore.close
    @ (s: store_state) -> result[void, string]
    + closes the underlying database handle
    # lifecycle
