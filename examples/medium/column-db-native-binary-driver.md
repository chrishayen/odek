# Requirement: "a client driver with a native binary interface for a columnar analytics database"

Speaks the database's binary protocol over tcp, runs blocks of inserts and queries. Raw sockets live in std.

std
  std.net
    std.net.tcp_connect
      fn (host: string, port: i32) -> result[tcp_conn, string]
      + opens a tcp connection
      - returns error on connect failure
      # network
    std.net.tcp_send
      fn (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes bytes to a connection
      - returns error on broken pipe
      # network
    std.net.tcp_recv
      fn (conn: tcp_conn, n: i32) -> result[bytes, string]
      + reads up to n bytes
      - returns error on connection close
      # network

column_db
  column_db.connect
    fn (host: string, port: i32, database: string, user: string, password: string) -> result[column_session, string]
    + performs the protocol handshake and returns a ready session
    - returns error on credentials or network failure
    # session
    -> std.net.tcp_connect
    -> std.net.tcp_send
    -> std.net.tcp_recv
  column_db.close
    fn (session: column_session) -> void
    + releases the underlying connection
    # session
  column_db.query
    fn (session: column_session, sql: string) -> result[list[map[string, string]], string]
    + runs a select and returns rows as name-to-value maps
    - returns error on server error
    # query
    -> std.net.tcp_send
    -> std.net.tcp_recv
  column_db.insert_block
    fn (session: column_session, table: string, columns: list[string], rows: list[list[string]]) -> result[i64, string]
    + sends a columnar block insert and returns the number of rows written
    - returns error when column count does not match row width
    # writes
    -> std.net.tcp_send
    -> std.net.tcp_recv
  column_db.execute
    fn (session: column_session, sql: string) -> result[void, string]
    + runs a non-query statement such as DDL
    - returns error on server error
    # ddl
    -> std.net.tcp_send
    -> std.net.tcp_recv
