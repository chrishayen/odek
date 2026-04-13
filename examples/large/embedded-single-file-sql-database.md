# Requirement: "an embedded single-file relational database with a SQL subset"

One file on disk holds all pages. A tiny SQL parser compiles statements into a plan that runs against a paged B-tree.

std
  std.fs
    std.fs.open_rw
      @ (path: string) -> result[file_handle, string]
      + opens or creates a file for read and write
      # filesystem
    std.fs.read_at
      @ (handle: file_handle, offset: i64, n: i32) -> result[bytes, string]
      + reads n bytes at the given offset
      # filesystem
    std.fs.write_at
      @ (handle: file_handle, offset: i64, data: bytes) -> result[void, string]
      + writes bytes at the given offset
      # filesystem
    std.fs.fsync
      @ (handle: file_handle) -> result[void, string]
      + flushes pending writes to durable storage
      # filesystem
  std.hash
    std.hash.crc32
      @ (data: bytes) -> u32
      + computes a CRC32 checksum
      # hashing

sqlitey
  sqlitey.open
    @ (path: string) -> result[db_state, string]
    + opens the database file, creating it with a header page when new
    - returns error when the file header is not recognized
    # lifecycle
    -> std.fs.open_rw
    -> std.fs.read_at
    -> std.fs.write_at
  sqlitey.close
    @ (state: db_state) -> result[void, string]
    + flushes the page cache and releases the file handle
    # lifecycle
    -> std.fs.fsync
  sqlitey.read_page
    @ (state: db_state, page_no: i32) -> result[tuple[bytes, db_state], string]
    + returns the bytes of the given page, caching it on miss
    # paging
    -> std.fs.read_at
  sqlitey.write_page
    @ (state: db_state, page_no: i32, data: bytes) -> result[db_state, string]
    + writes a page and marks it dirty in the cache
    # paging
    -> std.fs.write_at
    -> std.hash.crc32
  sqlitey.tokenize_sql
    @ (source: string) -> result[list[sql_token], string]
    + recognizes keywords, identifiers, literals, and punctuation
    - returns error on unterminated string literals
    # parser
  sqlitey.parse_statement
    @ (tokens: list[sql_token]) -> result[sql_statement, string]
    + builds an AST for SELECT, INSERT, UPDATE, DELETE, CREATE TABLE
    - returns error on unsupported syntax
    # parser
  sqlitey.plan_statement
    @ (state: db_state, stmt: sql_statement) -> result[sql_plan, string]
    + resolves tables and columns and produces an executable plan
    - returns error when a referenced table does not exist
    # planner
  sqlitey.execute_plan
    @ (state: db_state, plan: sql_plan) -> result[tuple[sql_result, db_state], string]
    + runs the plan, returning rows for queries and affected-count for mutations
    - returns error on type mismatches between literal and column
    # executor
  sqlitey.exec
    @ (state: db_state, sql: string) -> result[tuple[sql_result, db_state], string]
    + tokenizes, parses, plans, and executes a SQL statement
    # convenience
  sqlitey.btree_insert
    @ (state: db_state, root_page: i32, key: bytes, value: bytes) -> result[db_state, string]
    + inserts a key-value pair, splitting pages as needed
    # storage
  sqlitey.btree_lookup
    @ (state: db_state, root_page: i32, key: bytes) -> optional[bytes]
    + returns the value for the given key
    # storage
  sqlitey.btree_delete
    @ (state: db_state, root_page: i32, key: bytes) -> result[db_state, string]
    + removes the key, rebalancing underfull pages
    # storage
  sqlitey.begin_tx
    @ (state: db_state) -> result[db_state, string]
    + opens an implicit write transaction with page shadowing
    # transactions
  sqlitey.commit_tx
    @ (state: db_state) -> result[db_state, string]
    + atomically replaces shadow pages and fsyncs
    # transactions
    -> std.fs.fsync
  sqlitey.rollback_tx
    @ (state: db_state) -> db_state
    + discards shadow pages, leaving the committed state intact
    # transactions
