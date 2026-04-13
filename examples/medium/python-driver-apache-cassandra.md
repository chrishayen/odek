# Requirement: "a client driver for a distributed wide-column database"

Connects to cluster nodes, prepares statements, runs CQL queries. The actual wire transport is a std primitive so tests can substitute a fake node.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: i32) -> result[tcp_conn, string]
      + opens a tcp connection
      - returns error on dns or connect failure
      # network
    std.net.tcp_send
      @ (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes bytes to a connection
      - returns error on broken pipe
      # network
    std.net.tcp_recv
      @ (conn: tcp_conn, n: i32) -> result[bytes, string]
      + reads up to n bytes
      - returns error on connection close
      # network

wide_column
  wide_column.new_cluster
    @ (contact_points: list[string], port: i32) -> cluster_state
    + records contact points for lazy connection
    # cluster
  wide_column.connect
    @ (cluster: cluster_state, keyspace: string) -> result[session_state, string]
    + establishes a session bound to a keyspace
    - returns error when no contact point is reachable
    # session
    -> std.net.tcp_connect
  wide_column.close
    @ (session: session_state) -> void
    + closes all open connections in the session
    # session
  wide_column.prepare
    @ (session: session_state, cql: string) -> result[prepared_statement, string]
    + returns a prepared statement handle
    - returns error when the server rejects the cql
    # statements
    -> std.net.tcp_send
    -> std.net.tcp_recv
  wide_column.execute
    @ (session: session_state, prepared: prepared_statement, params: list[string]) -> result[list[map[string, string]], string]
    + binds params and executes the statement, returning rows as maps
    - returns error on server error or parameter arity mismatch
    # statements
    -> std.net.tcp_send
    -> std.net.tcp_recv
  wide_column.execute_simple
    @ (session: session_state, cql: string) -> result[list[map[string, string]], string]
    + executes a one-shot query with no parameters
    - returns error on server error
    # statements
    -> std.net.tcp_send
    -> std.net.tcp_recv
