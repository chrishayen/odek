# Requirement: "a native client and SQL driver for a distributed database"

Connect to a cluster, maintain a session pool, run typed queries and parameterized statements, and expose a generic SQL cursor adapter.

std
  std.net
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[conn_state, string]
      + opens a TCP connection to host:port
      - returns error on dns, timeout, or refused
      # networking
    std.net.read
      fn (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      - returns error when the connection is closed
      # networking
    std.net.write
      fn (conn: conn_state, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      - returns error when the connection is closed
      # networking
    std.net.close
      fn (conn: conn_state) -> void
      + closes the connection
      # networking
  std.encoding
    std.encoding.varint_encode
      fn (value: u64) -> bytes
      + encodes an unsigned integer as little-endian base-128 varint
      # encoding
    std.encoding.varint_decode
      fn (data: bytes) -> result[tuple[u64, i32], string]
      + decodes a varint and returns the value and number of bytes consumed
      - returns error on truncated or oversized input
      # encoding
  std.sync
    std.sync.mutex_new
      fn () -> mutex_state
      + creates an unlocked mutex
      # concurrency
    std.sync.mutex_lock
      fn (m: mutex_state) -> void
      + acquires the mutex, blocking until available
      # concurrency
    std.sync.mutex_unlock
      fn (m: mutex_state) -> void
      + releases the mutex
      # concurrency

ydb_driver
  ydb_driver.connect
    fn (endpoints: list[string], database: string) -> result[driver_state, string]
    + dials the first reachable endpoint and records the database name
    - returns error when no endpoint is reachable
    # connection
    -> std.net.dial_tcp
  ydb_driver.close
    fn (state: driver_state) -> void
    + closes every pooled session and the underlying connections
    # lifecycle
    -> std.net.close
  ydb_driver.session_pool_new
    fn (state: driver_state, max_size: i32) -> driver_state
    + configures an empty session pool with the given cap
    # pooling
    -> std.sync.mutex_new
  ydb_driver.acquire_session
    fn (state: driver_state) -> result[session_handle, string]
    + returns an idle session or creates one when below the cap
    - returns error when the pool is exhausted and the wait times out
    # pooling
    -> std.sync.mutex_lock
    -> std.sync.mutex_unlock
  ydb_driver.release_session
    fn (state: driver_state, session: session_handle) -> void
    + marks the session as idle and returns it to the pool
    # pooling
    -> std.sync.mutex_lock
    -> std.sync.mutex_unlock
  ydb_driver.encode_query
    fn (sql: string, params: map[string, typed_value]) -> bytes
    + builds the length-prefixed request frame for a parameterized query
    # protocol
    -> std.encoding.varint_encode
  ydb_driver.decode_result_set
    fn (frame: bytes) -> result[result_set, string]
    + parses a response frame into a typed result set
    - returns error on truncated frames or unknown column types
    # protocol
    -> std.encoding.varint_decode
  ydb_driver.execute
    fn (state: driver_state, session: session_handle, sql: string, params: map[string, typed_value]) -> result[result_set, string]
    + sends the query on the session and returns the decoded result set
    - returns error when the request or response fails
    # execution
    -> std.net.write
    -> std.net.read
  ydb_driver.begin_tx
    fn (state: driver_state, session: session_handle, mode: string) -> result[tx_handle, string]
    + opens a transaction in the given isolation mode
    - returns error when mode is unknown
    # transactions
    -> std.net.write
    -> std.net.read
  ydb_driver.commit_tx
    fn (state: driver_state, tx: tx_handle) -> result[void, string]
    + commits the transaction
    - returns error when the server reports conflict
    # transactions
    -> std.net.write
    -> std.net.read
  ydb_driver.rollback_tx
    fn (state: driver_state, tx: tx_handle) -> result[void, string]
    + rolls back the transaction
    # transactions
    -> std.net.write
    -> std.net.read
  ydb_driver.sql_open
    fn (endpoints: list[string], database: string) -> result[sql_conn, string]
    + returns a generic SQL adapter on top of the native driver
    # sql_adapter
  ydb_driver.sql_query
    fn (conn: sql_conn, query: string, args: list[typed_value]) -> result[sql_rows, string]
    + runs a query and returns a cursor over the result rows
    - returns error on protocol or execution failure
    # sql_adapter
