# Requirement: "an in-process SQL database library"

A small relational engine: a B-tree-backed page store, a SQL parser, a query planner, and an executor. Parsing and planning are pure; storage is behind an injected pager.

std
  std.encoding
    std.encoding.varint_encode
      fn (value: i64) -> bytes
      + returns the little-endian varint encoding of value
      # encoding
    std.encoding.varint_decode
      fn (data: bytes, offset: i64) -> result[tuple[i64, i64], string]
      + returns (value, new_offset)
      - returns error on truncated input
      # encoding
  std.hash
    std.hash.fnv1a_64
      fn (data: bytes) -> u64
      + returns 64-bit FNV-1a hash
      # hashing

sqldb
  sqldb.new
    fn (pager: pager_handle) -> db_state
    + creates an empty database bound to the given pager
    # construction
  sqldb.tokenize
    fn (sql: string) -> result[list[sql_token], string]
    + returns the token stream for a SQL statement
    - returns error on an unterminated string literal
    # lexing
  sqldb.parse
    fn (tokens: list[sql_token]) -> result[sql_statement, string]
    + returns an AST for SELECT, INSERT, UPDATE, DELETE, CREATE TABLE, DROP TABLE
    - returns error on a syntax error with position
    # parsing
  sqldb.plan
    fn (state: db_state, stmt: sql_statement) -> result[query_plan, string]
    + returns a plan: seq-scan, index-scan, or index-lookup with filters and projections
    - returns error when a referenced table or column does not exist
    # planning
  sqldb.execute
    fn (state: db_state, plan: query_plan) -> result[tuple[query_result, db_state], string]
    + runs the plan and returns rows and the updated db state
    - returns error on a unique-constraint violation
    # execution
  sqldb.create_table
    fn (state: db_state, name: string, columns: list[column_def]) -> result[db_state, string]
    + registers a new table in the catalog and allocates its root page
    - returns error when a table with the name already exists
    # ddl
  sqldb.btree_insert
    fn (state: db_state, root_page: i64, key: bytes, value: bytes) -> result[db_state, string]
    + inserts a key/value into the B-tree, splitting leaves as needed
    # storage
    -> std.encoding.varint_encode
  sqldb.btree_lookup
    fn (state: db_state, root_page: i64, key: bytes) -> optional[bytes]
    + returns the value bound to key if present
    # storage
    -> std.encoding.varint_decode
  sqldb.btree_range
    fn (state: db_state, root_page: i64, start: bytes, end: bytes) -> list[tuple[bytes, bytes]]
    + returns all key/value pairs in [start, end] in key order
    # storage
  sqldb.begin_transaction
    fn (state: db_state) -> db_state
    + opens a new transaction; subsequent writes are staged
    # transaction
  sqldb.commit
    fn (state: db_state) -> result[db_state, string]
    + flushes staged writes through the pager
    - returns error when no transaction is open
    # transaction
  sqldb.rollback
    fn (state: db_state) -> db_state
    + discards staged writes
    # transaction
  sqldb.row_hash
    fn (row: list[bytes]) -> u64
    + returns a stable hash used for unique-index dedup
    # indexing
    -> std.hash.fnv1a_64
