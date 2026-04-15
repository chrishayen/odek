# Requirement: "a client driver for a wide-column distributed database"

A connection-oriented driver: cluster configuration, session management, prepared statements, query execution, and result iteration. Wire framing lives behind pluggable transport; this library handles protocol framing, routing, and result decoding.

std
  std.net
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[conn_handle, string]
      + opens a TCP connection to host:port
      - returns error when the connection cannot be established
      # networking
    std.net.write_all
      fn (conn: conn_handle, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      - returns error on connection failure
      # networking
    std.net.read_exact
      fn (conn: conn_handle, n: i32) -> result[bytes, string]
      + reads exactly n bytes from the connection
      - returns error on short read or connection failure
      # networking
  std.encoding
    std.encoding.read_u32_be
      fn (data: bytes, offset: i32) -> result[u32, string]
      + reads a big-endian u32 at offset
      - returns error when offset is out of range
      # encoding
    std.encoding.write_u32_be
      fn (value: u32) -> bytes
      + serializes a u32 as 4 big-endian bytes
      # encoding

cassandra_driver
  cassandra_driver.new_cluster
    fn (contact_points: list[string], port: i32) -> cluster_state
    + creates a cluster configuration with the given contact points and port
    # construction
  cassandra_driver.connect
    fn (cluster: cluster_state) -> result[session_state, string]
    + establishes connections to the contact points and returns a session
    - returns error when every contact point is unreachable
    # session
    -> std.net.dial_tcp
  cassandra_driver.startup
    fn (session: session_state) -> result[session_state, string]
    + performs the protocol startup handshake
    - returns error when the server rejects the startup frame
    # protocol
    -> std.net.write_all
    -> std.net.read_exact
  cassandra_driver.prepare
    fn (session: session_state, statement: string) -> result[prepared_id, string]
    + prepares a query and returns its opaque handle
    - returns error when the server reports a syntax error
    # statement
    -> std.net.write_all
    -> std.net.read_exact
  cassandra_driver.execute
    fn (session: session_state, prepared: prepared_id, params: list[bytes]) -> result[result_set, string]
    + executes a prepared statement with bound parameters
    - returns error when the server responds with an error frame
    # query
    -> std.net.write_all
    -> std.net.read_exact
    -> std.encoding.write_u32_be
  cassandra_driver.query
    fn (session: session_state, cql: string) -> result[result_set, string]
    + sends an unprepared query string and returns the result set
    - returns error on protocol or server-side failures
    # query
    -> std.net.write_all
    -> std.net.read_exact
  cassandra_driver.row_count
    fn (results: result_set) -> i32
    + returns the number of rows in a result set
    # results
  cassandra_driver.get_row
    fn (results: result_set, index: i32) -> result[list[bytes], string]
    + returns the raw column values of a row
    - returns error when index is out of range
    # results
  cassandra_driver.encode_frame
    fn (opcode: u8, stream_id: i16, body: bytes) -> bytes
    + builds a protocol frame from opcode, stream id, and body
    # protocol
    -> std.encoding.write_u32_be
  cassandra_driver.decode_frame
    fn (data: bytes) -> result[tuple[u8, i16, bytes], string]
    + parses a protocol frame into (opcode, stream_id, body)
    - returns error on truncated or malformed frames
    # protocol
    -> std.encoding.read_u32_be
  cassandra_driver.close
    fn (session: session_state) -> result[void, string]
    + tears down all open connections in the session
    # session
